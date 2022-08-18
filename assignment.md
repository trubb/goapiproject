# Code test Golang - devies

## Mission

Create a RESTful web API in Go that connects to an underlying data layer of your choice (SQL, NoSQL etc).
This service should be able to get and set data to and from the database.

### Requirements

1. You need to expose various endpoints with REST-based routing.
2. Endpoints will have to be accessed and tested either through Postman or any other front-end client.
3. Your API should follow a common/best-practice pattern, or any other modern design pattern.
4. You are free to choose the theme of the service. You may find an example on the next page.

### Purpose

The purpose of this test is to find out how well versed you are in Golang, backend development and its best practices. We send a test of this magnitude to judge your ability to take home a task, do your research on it, and deliver the results — in your own pace.

### Evaluation

What we will judge is:

- Modernity of code
- Following best practices
- How well you can explain (and document) your code.
  - Remember; it’s important to know why something is written the way it is, rather than how.
- Feature-richness.
  - Does your code provide error handling?
  - Input validation?

## Example of API themes

Build a product inventory system, with entities such as Customer, Product and
Orders. You should be able to put some orders containing certain parameters, like
list of product, etc.

Customer
Money - get / set
Orders - get / set
... etc ...

->

Order
Products - get / set
Product Price - get
... etc ....

->

Product
Price - get / set
Product Details - get / set
... etc ....

You are free to design this the way you want, and these examples serve as a help
for you to get started.
