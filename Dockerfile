FROM golang
COPY . /app
RUN cd /app && go build -o run main.go
