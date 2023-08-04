# Standard Library/Golang: Starter API Code Sample

This Golang code sample demonstrates how to build an API server using Standard Library that is secure by design.

Visit the ["Standard Library/Golang Code Samples: API Security in Action"](https://developer.auth0.com/resources/code-samples/api/standard-library) section of the ["Auth0 Developer Resources"](https://developer.auth0.com/resources) to explore how you can secure Standard Library applications written in Golang by implementing endpoint protection and authorization with Auth0.

## Why Use Auth0?

Auth0 is a flexible drop-in solution to add authentication and authorization services to your applications. Your team and organization can avoid the cost, time, and risk that come with building your own solution to authenticate and authorize users. We offer tons of guidance and SDKs for you to get started and [integrate Auth0 into your stack easily](https://developer.auth0.com/resources/code-samples/full-stack).

## Set Up and Run the Standard Library Project

Download module dependencies to local cache:

```bash
go mod download
```

Create a `.env` file under the root project directory and populate it with the following environment variables:

```bash
PORT=6060
CLIENT_ORIGIN_URL=http://localhost:4040
```

Finally, execute this command to start the api server:

```bash
go run ./cmd/api-server/
```
