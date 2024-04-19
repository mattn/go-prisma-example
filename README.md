# go-prisma-example

Example Go app using [prisma](https://www.prisma.io/), [lit-html](https://lit-html.polymer-project.org/), and [ky](https://github.com/sindresorhus/ky).

## Usage

First of all, install PostgreSQL.
Then try following commands.

```console
$ go generate
$ go build
$ echo 'DATABASE_URL=postgresql://yourname:yourname@localhost:5432/postgres' > .env
$ npx prisma db push
$ ./go-prisma-example
```

Then you can use [TODO app](http://localhost:8989)! :tada:

## Reference

* [Blog post(ja)](https://zenn.dev/mattn/articles/1c4eb193d81a3a)

## License

MIT

assets/app.js is based on @ryohey's [lit-html-todo](https://github.com/ryohey/lit-html-todo)

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
