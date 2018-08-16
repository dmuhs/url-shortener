build:
	cd url-shortener && go build -i .

run:
	cd url-shortener && go build -i . && ./url-shortener

clean:
	rm -f url-shortener/url-shortener
	rm -f url-shortener/urls.db
