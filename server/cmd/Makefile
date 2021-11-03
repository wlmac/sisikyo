OAUTH_CLIENT_ID="XNc77hFC2VjQJUqW0vLrlsH6UInm2Uivf6u52GE3"
OAUTH_CLIENT_SECRET="250WW680mvzDZuR2I1gGncM3Ezlw85DHUA5SjCgFcidLCvy8qh5VTx96cxJSIjlLXMomfLkWvShuuPx1ICVfqVfLfEyCJhbkR05lf4BXkX0ckzbMkYzt6ujNvwsbDXLj"

run-local: clean server
	./server \
		--api-url http://localhost:8000/api/ \
		--o-url http://localhost:8000/o/ \
		--db-driver sqlite3 \
		--db-source test.db \
		--o-id ${OAUTH_CLIENT_ID} \
		--o-secret ${OAUTH_CLIENT_SECRET}

clean:
	rm -f ./server

server:
	go build -o=server --tags=sqlite ./main.go

.PHONY: run-local clean;
