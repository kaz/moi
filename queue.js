"use strict";

const Redis = require("ioredis");

(async () => {
	const redis = new Redis();

	for(let i = 1111111; i <= 9999999; i++){
		const num = `${i}`.split("").join(" ? ");
		if(num.includes("0")){
			continue;
		}
		await redis.lpush("que", num);
	}

	redis.disconnect();
})();
