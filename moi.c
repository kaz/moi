#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <arpa/inet.h>
#include <netinet/in.h>
#include <openssl/ssl.h>

#ifdef DEBUG
#define LOG(io, data, len) printf("\n\n---------- %s ----------\n\n", (io) ? "RECV" : "SEND"), write(1, (data), (len))
#else
#define LOG(io, data, len)
#endif

const char* handshake = "\
GET /internships/2019/games?level=3 HTTP/1.1\n\
Host: apiv2.twitcasting.tv\n\
Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImp0aSI6ImI1NTc5YWE2N2VlNDI3MDBmODhiNTQ5ZjVjNjA2MTgxYTM1ODFhZTEzODE1MmM1OTZjMDNjNTJlOTQ4ZjUyOTkyMzRiYjFiN2JjZmE3ZWEzIn0.eyJhdWQiOiIxODIyMjQ5MzguMjNhNzJmNDA2NzI4M2I0OWY5NjZmOTMyMzViMTg2NDQzN2VjNWY2YTlmY2M5NjVlOGIzOTM5MGRmNWQ2YWE5NCIsImp0aSI6ImI1NTc5YWE2N2VlNDI3MDBmODhiNTQ5ZjVjNjA2MTgxYTM1ODFhZTEzODE1MmM1OTZjMDNjNTJlOTQ4ZjUyOTkyMzRiYjFiN2JjZmE3ZWEzIiwiaWF0IjoxNTYyMzMyNzcyLCJuYmYiOjE1NjIzMzI3NzIsImV4cCI6MTU3Nzg4NDc3MSwic3ViIjoiYzphbGljZV9nIiwic2NvcGVzIjpbInJlYWQiXX0.GoWeDo-ZswQ1sp_ejnj9rKT2MwPaZwYpqC_w9GS5r1_bJ-aiPPd8rwUPUY3VFphkSVkVpZRwyGq3bxc1Rx1CQMh5_xBavoaMtCr1iik4YZPTIJJJBjXfpgQTOqRGgKA5HEvz84d4XnQGeFNBCw7zbpfOiBmiByHDfulh_SjkI-AwCUPJ4vaXHVjcHHtqlQfZ5jxYwLH2Zv0duwlsDMxR-tWU70TGeFV71yByE53fL4s6Heg607BeFDFhIvNmkMOULiv5xOrnJlmGrxySfflZn4KbStRydypvfgAc2Kbkz1YUQatQazN4hvCrr-otpyPccdKBLF8cf4UVpwSpCSClzQ\n\
\n";

int prep_buf_len;
char* prep_buf;
SSL* prep_ssl;

SSL* open_connection() {
	int sockfd = socket(AF_INET, SOCK_STREAM, 0);

	struct sockaddr_in addr;
	addr.sin_family = AF_INET;
	addr.sin_port = htons(443);
	inet_aton("202.239.41.35", &addr.sin_addr);
	connect(sockfd, (const struct sockaddr*) &addr, 32);

	SSL_library_init();

	SSL_CTX* ctx = SSL_CTX_new(SSLv23_client_method());
	SSL_CTX_set_cipher_list(ctx, "AES128-GCM-SHA256");

	SSL* ssl = SSL_new(ctx);
	SSL_set_fd(ssl, sockfd);
	SSL_connect(ssl);

	return ssl;
}

char* start_game() {
	int read;
	int buf_len = 2048;
	char *buf = malloc(buf_len);

	SSL* ssl = open_connection();
	SSL_write(ssl, handshake, strlen(handshake));

	read = SSL_read(ssl, buf, buf_len);
	read = SSL_read(ssl, buf, buf_len);

	buf[read] = 0;
	return buf;
}

void prepare() {
	prep_buf_len = 2048;
	prep_buf = malloc(prep_buf_len);
	prep_ssl = open_connection();
}
char* answer(char *data) {
	int read;
	int buf_len = prep_buf_len;
	char* buf = prep_buf;
	SSL* ssl = prep_ssl;

	LOG(0, data, strlen(data));
	SSL_write(ssl, data, strlen(data));

	read = SSL_read(ssl, buf, buf_len);
	LOG(1, buf, read);

	read = SSL_read(ssl, buf, buf_len);
	LOG(1, buf, read);

	buf[read] = 0;
	return buf;
}
