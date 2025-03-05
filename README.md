# Cart Service Assessment

This repo contains the code to solve the cart service assessment.

## Considerations

The repo goal aims to show some ideas and possible practice that can be used to expose a RESTful API with a go stack

In the code, some comments about the decision made are present. Moreover, here a list of guideline followed during this specific assessment:
- For simplicity, no specific http handling library was used (like Fiber or Gin). The standard Go library was enough for the scope of the project
- The DB of choice is SqlLite3. Being a cart service having a relational database was nice in case of future need where more strict transaction handling is needed. Nothing was implemented in the assessment since not strictly requested
- The entities modelled are very small, hanging only the needed field for the API needs. Field like *name* of a product, or table to store the order detail was not implemented by choice since was adding nothing to the "decision-making" around the assessment
- Were possible, the test used real database (in this case, temporary SqlLite3 one) since testing against real implementations instead of mocking brings to more valuable test result. The test are a mix of proper unit test, some integration one and even a couple of e2e, to test the whole traversing of information in the system.
- Resources like database migration definition and the main `Dockerfile` are purposely built for this assessment, not for a real case scenario.

## How to run

Following the requested path, here the command to run the application, leveraging Docker:

```sh
docker build -t assessment-cart-service
docker run -v $(pwd):/mnt -p 9090:9090 -w /mnt assessment-cart-service:latest ./scripts/build.sh
docker run -v $(pwd):/mnt -p 9090:9090 -w /mnt assessment-cart-service:latest ./scripts/tests.sh
docker run -v $(pwd):/mnt -p 9090:9090 -w /mnt assessment-cart-service:latest ./scripts/run.sh
```