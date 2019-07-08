"use strict";

const cluster = require("cluster");
const Redis = require("ioredis");

const ops = ["+", "-", "*", "/"];
let pats = [""];

for(let i = 0; i < 6; i++){
	const npats = [];
	pats.forEach(pat => {
		ops.forEach(op => {
			npats.push(pat+op);
		});
	});
	pats = npats;
}

if(cluster.isMaster){
	require("os").cpus().forEach(() => {
		cluster.fork();
	});

	console.log(`Master ${process.pid} started`);
}else{
	(async () => {
		const redis = new Redis();

		while(true){
			const num = await redis.lpop("que");
			const expr = num.split("");

			const vals = [];
			pats.forEach(pat => {
				pat.split("").forEach((c, i) => {
					expr[2+4*i] = c;
				});

				const val = eval(expr.join(""));
				if(val - parseInt(val) != 0){
					return;
				}

				vals.push(val, pat);
			});

			await redis.hmset(num, ...vals);
		}
	})();

	console.log(`Worker ${process.pid} started`);
}
