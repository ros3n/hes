# HES

## Problem

Creation of a web service that provides users with fault-tolerant mailing services. The solution should provide an abstraction over multiple service providers. In case one of the providers goes down, a failover mechanism should utilise another one in order to send the emails. The service should be implemented in a way that makes it possible to be used by other services, i.e. by implementing a REST API. The solution should utilise some form of authentication. Also, it should be a cloud-native solution that could be deployed easily.

## Solution

The resulting software consists of two services:
- a user-facing REST API that makes it possible to create emails and schedule sends (documented [here](https://app.swaggerhub.com/apis/rafalr/HES/1.0.0))
- an internal mailer service that accepts send requests and uses available providers to deliver messages

#### API

The API is implemented using the Go programming language. The language has been chosen because of its robustness when it comes to parallel processing and the developer's experience with it (about 1.5 years, ~50% of work at the previous company).

The API provides REST endpoints that make it possible to create email resources, list them, and schedule particular emails for sends (the number of CRUD operations is limited to essentials due to the tight deadline). The data must be provided in the JSON format. In order to be stored, the data must pass a validation process. The API stores the data in an SQL database. When a user requests an email send, the API communicates that information to the mailer service using a gRPC call and marks the email as queued for send. The API provides an RPC callback method that is used by the mailer service to report send statuses.

Access to the API is guarded by an authentication module. The current implementation utilises BasicAuth, but the middleware has been created in a way that makes it easy to integrate with an external IAM service.

#### Mailer

The mailer is implemented using the Go programming language as well. It has been chosen for the same reasons.

The mailer provides an RPC method that makes it possible for the API to request email sends. Communication is done using gRPC. After receiving a send request the mailer's manager creates a new worker and tasks it to send the email. The worker chooses a service provider using round-robin strategy and tries to send the email. In case of failure it fail-overs to the next provider. The current implementation supports SendGrid and MailGun, but additional providers could be added easily. After a successful send (or a failure after retry) workers report send statuses to the manager, which in turn informs the API about the result.

The current implementation should be considered an MVP, due to lack of durable jobs. The service providers do not support any form of request ids, hence the job of ensuring idempotence is the client's responsibility. That would require development of a persistence layer for the jobs, a garbage collection system, and a synchronisation mechanism. Implementing such solution given the deadline could be a too long shot, hence the developer decided to develop an MVP version.

##### Scalability
Because the application is split into two services and only one of them relies on a database, the system should be fairly scalable. It is possible to scale the API and the mailer independently, for example by adjusting the number of Kubernetes pods created for each service.

##### The code
Apart from the code stored in the *vendor/* directory and the *lib/communication/communication.pb.go* file, all of the go code has been written by the developer.

#### Deployment

The solution provides Kubernetes configuration files that can be used to deploy the system. They live in _k8s/_ directory. In order to make the deployments work one must first create a secrets resource called _hes-secrets_ (see k8s/secrets/hes-secrets.yml.example). Since the MVP does not validate presence of the required configuration variables, please make sure that the secret contains all of them.

There is a Makefile that automates image building and pushing it to the DockerHub. It provides the following tasks:
- make build-api
- make build-mailer
- make dockerize-api
- make dockerize-mailer
- make push-api
- make push-mailer

The Makefile expects *DOCKER_USER* env variable to be present.

Database migrations are handled by ActiveRecord. There is a Rakefile that provides the following tasks:
- bundle exec rake db:create
- bundle exec rake db:migrate
- bundle exec rake db:rollback

Before running migrations one must first setup bundler: gem install bundler && bundle install. The Rakefile expects *HES_DATABASE_URL* to be set.

#### Staging
There is a staging Kubernetes claster available. It is deployed on DigitalOcean. For further information about how to use it see the [docs](https://app.swaggerhub.com/apis/rafalr/HES/1.0.0).
