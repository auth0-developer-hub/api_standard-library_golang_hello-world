# Hello World API: Golang #

This repository contains a Golang API server. You'll secure these APIs with Auth0 to practice making secure API calls
from a client application.

## Get Started

### Set up the project

Install the project dependencies:

```bash
go mod download
```

Create the `env.yaml` file at the root of the project and populate it with the values of `auth0-audience` and
`auth0-domain` parameters as in the following example.

```yaml
auth0-audience: my_auth0_api_audience
auth0-domain: my_auth0_tenant_domain
```

The following sections further details the steps to identify the values for these variables within your setup.

### Register a Golang API with Auth0

- Open the [APIs](https://manage.auth0.com/#/apis) section of the Auth0 Dashboard.

- Click on the **Create API** button.

- Provide a **Name** value such as _Hello World API Server_.

- Set its **Identifier** to `https://api.example.com` or any other value of your liking.

- Leave the signing algorithm as `RS256` as it's the best option from a security standpoint.

- Click on the **Create** button.

> View ["Register APIs" document](https://auth0.com/docs/get-started/set-up-apis) for more details.

### Get API configuration values

Head back to your Auth0 API page, and follow these steps to get the Auth0 Audience:

![Get the Auth0 Audience to configure an API](https://images.ctfassets.net/23aumh6u8s0i/1CaZWZK062axeF2cpr884K/cbf29676284e12f8e234545de05dac58/get-the-auth0-audience)

- Click on the "Settings" tab.

- Locate the "Identifier" field and copy its value.

- Paste the Auth0 domain value as the value of `auth0-audience` in `env.yaml`.

Now, follow these steps to get the Auth0 Domain value:

![Get the Auth0 Domain to configure an API](https://images.ctfassets.net/23aumh6u8s0i/37J4EUXKJWZxHIyxAQ8SYI/d968d967b5e954fc400163638ac2625f/get-the-auth0-domain)

- Click on the "Test" tab.

- Locate the section called "Asking Auth0 for tokens from my application".

- Click on the cURL tab to show a mock `POST` request.

- Locate your Auth0 domain, which is part of the `--url` parameter value: `tenant-name.region.auth0.com`.

- Paste the Auth0 domain value as the value of `auth0-domain` in `env.yaml`.

**Tips to get the Auth0 Domain**

- The Auth0 Domain is the substring between the protocol, `https://` and the path `/oauth/token`.

- The Auth0 Domain follows this pattern: `tenant-name.region.auth0.com`.

- The `region` subdomain (`au`, `us`, or `eu`) is optional. Some Auth0 Domains don't have it.

### Run the API server:

Run the API server by using any of the following commands. Please replace `my_auth0_api_audience`
and `my_auth0_api_domain` values appropriate as per your setup using instructions above.

```shell
# with a populated env.yaml/dev.yaml present in current directory

./run.sh

```

# OR, you may also expose the values as environment variable upfront
export AUTH0_AUDIENCE=my_auth0_api_audience
export AUTH0_DOMAIN=my_auth0_api_domain
go run .
```

In case the auth0 `audience` and `domain` values are present at multiple sources, the following precedence order is used
to determine effective values:

- CLI parameters values has the highest precedence

- Shell/Environment variables takes the next precedence

- Values from `env.yaml` has the lowest precedence

## Test the Protected Endpoints

You can get an access token from the Auth0 Dashboard to test making a secure call to your protected API endpoints.

Head back to your Auth0 API page and click on the "Test" tab.

Locate the section called "Sending the token to the API".

Click on the cURL tab of the code box.

Copy the sample cURL command:

```bash
curl --request GET \
  --url http://path_to_your_api/ \
  --header 'authorization: Bearer really-long-string-which-is-test-your-access-token'
```

Replace the value of `http://path_to_your_api/` with your protected API endpoint path (you can find all the available
API endpoints in the next section) and execute the command. You should receive back a successful response from the
server.

You can try out any of our full stack demos to see the client-server Auth0 workflow in action using your preferred
front-end and back-end technologies.

## Test the Admin Endpoint

The `/admin` endpoint requires the access token to contain the `read:admin-messages` permission. The best way to
simulate this client-server secured request is to use any of the Hello World client demo apps to log in as a user that
has that permission.

Use following steps to create a user with `read:admin-messages` permission.

- Head back to your Auth0 API page, click on "Settings" tab. Toggle "Enable RBAC" and "Add Permissions in the Access
  Token" to enabled state. Click "Save" at bottom of this page.

- Click "Permissions" tab on your Auth0 API page. Add a new permission with `Scope` as `read:admin-messages`, use
  any `Description` of your choice.

- Create an `admin` role from the Auth0 dashboard, under "User Management"

- Go to "Permissions" tab of the `admin` role and click "Add Permissions". Select your API from the dropdown and add
  the `read:admin-messages` permission to the role.

- Assign the `admin` role to any existing user. You can also create a new user and assign `admin` role to the new user
  instead.

## API Endpoints

### 🔓 Get public message

```shell
GET /api/messages/public
```

#### Response

```shell
Status: 200 OK
```

```json
{
  "message": "The API doesn't require an access token to share this message."
}
```

### 🔓 Get protected message

> You need to protect this endpoint using Auth0.

```shell
GET /api/messages/protected
```

#### Response

```shell
Status: 200 OK
```

```json
{
  "message": "The API successfully validated your access token."
}
```

### 🔓 Get admin message

> You need to protect this endpoint using Auth0 and Role-Based Access Control (RBAC).

```shell
GET /api/messages/admin
```

#### Response

```shell
Status: 200 OK
```

```json
{
  "message": "The API successfully recognized you as an admin."
}
```

### Get Version

```shell
GET /version
```

#### Response

```shell
Status: 200 OK
```

```json
{"version":"7819552-dirty"}
```

### Get ping

```shell
GET /ping
```

#### Response

```shell
Status: 200 OK
```

```text
pong
```

## Data Storage

Data is stored in a postgres database by default.

## Query Logic

Find requests `GET /api/messages` uses a url query parser to allow very complex logic including AND, OR and precedence operators.
Ref: https://github.com/snowzach/queryp

## Swagger Documentation

Swagger documentation available at `http://<host>:<port>/api-docs/` path. The documentation is automatically generated.

## Profiler

Profiler is available at `http://<host>:6060/debug/pprof/` path. We can enable or disable using the config.

## Logs

The application logs can be found in `logs/oauth.log` file by default. We can change using the settings.go file.

Aditionally we can use Beats to ship logs from this file to Elastic search and use Kibana to visualize.

## Prometheus

The prometheus is available at `http://<host>:<port>/prometheus/` path. Now we can use Grafana to connect this datasource and visualize the metrics.

