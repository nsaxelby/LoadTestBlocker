# Instructions
The IAC folder provides a way to create a hosted version of the LoadTestBlocker in AWS (via an EC2 instance)

# Requirements
`terraform`

# Deployment
`terraform apply --auto-approve` is all you need

This creates a private key on your local device called `cert.pem`. This is your private key for your instance. You use it to connect to the ec2. Don't commit it to source control.

In the output of the terraform follow the link e.g.:

`web-connect = "http://54.170.83.204:8080"`

(It may take 20 seconds or so to startup)