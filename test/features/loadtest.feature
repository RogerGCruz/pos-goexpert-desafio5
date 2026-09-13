Feature: Load test CLI
  Validate the load test CLI produces correct report

  Scenario: successful load test against local server
    Given a test HTTP server that returns 200
    When I run the CLI with the server URL, requests 10 and concurrency 2
    Then the report shows total requests 10
    And the report shows 200 OK: 10

  Scenario: server returns 404
    Given a test HTTP server that returns 404
    When I run the CLI with the server URL, requests 5 and concurrency 2
    Then the report shows total requests 5
    And the report shows 404: 5

  Scenario: server returns 500
    Given a test HTTP server that returns 500
    When I run the CLI with the server URL, requests 5 and concurrency 2
    Then the report shows total requests 5
    And the report shows 500: 5

  Scenario: mixed status responses
    Given a test HTTP server that cycles through 200,404,500
    When I run the CLI with the server URL, requests 9 and concurrency 3
    Then the report shows total requests 9
    And the report shows 200: 3
    And the report shows 404: 3
    And the report shows 500: 3

  Scenario: server times out
    Given a test HTTP server that sleeps 20s before responding
    When I run the CLI with the server URL, requests 3 and concurrency 2
    Then the report shows total requests 3
    And the report shows 0: 3

  Scenario: latency impact
    Given a test HTTP server that sleeps 100ms before responding 200
    When I run the CLI with the server URL, requests 4 and concurrency 2
    Then the report shows total requests 4
    And the report total time is at least 200ms
