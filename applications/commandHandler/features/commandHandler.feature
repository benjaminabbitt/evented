Feature: command calls to business logic
  Scenario: Simple creation command
    Given a simple create command is added to the command queue
    When the command queue is sent to the coordinator
    Then the business logic should be called with the command queue