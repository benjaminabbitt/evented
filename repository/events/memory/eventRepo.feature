Feature: Event Repository
  The event repository stores events (things that have happened) related to a specific domain.
  It is keyed by the id of the object within the domain and the list of event sequence numbers that make up that object.  domain:id:sequence may be referred to as the coordinates of an event.

  Scenario Outline: Events with fixed coordinates should be storable and retrievable
    Given that we're working in the coordinates of a domain
    When I store the sample event:
      | id   | sequence   |
      | <id> | <sequence> |
    Then I should be able to retrieve it by its coordinates:
      | id   | sequence   |
      | <id> | <sequence> |
    Examples:
      | id                                   | sequence |
      | 7bf20643-77a7-48b7-b3bd-ac0c2649a8f3 | 0        |

#  Scenario: Events with forced sequence numbers should be stored and retrievable
#    Given that we're working in the coordinates of a domain
#    When I store the sample event:
#      | id                                   | sequence |
#      | 7bf20643-77a7-48b7-b3bd-ac0c2649a8f3 | force    |
#    Then I should be able to retrieve it by its coordinates:
#      | id                                   | sequence |
#      | 7bf20643-77a7-48b7-b3bd-ac0c2649a8f3 | 0        |
