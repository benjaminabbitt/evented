Feature: Create TODO
  In order to track my TODOs, I need to be able to create them to be tracked
  As a user
  I need to be able to create TODOs

  Scenario: Create a sample TODO
    Given a title of "Finish the TODO project"
    When I run this
    Then there should be an event created with the title "Finish the TODO project"
    And the sequence should be 0
    And the domain should be "todo"
    And the id should be set