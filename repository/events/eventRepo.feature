Feature: Event Repository
  The event repository stores events (things that have happened) related to a specific domain.
  It is keyed by the id of the object within the domain and the list of event sequence numbers that make up that object.  domain:id:sequence may be referred to as the coordinates of an event.

  The repository handles a single domain.

  Scenario Outline: Events with fixed coordinates should be storable and retrievable
    When I store the event:
      | id   | sequence   | time   |
      | <id> | <sequence> | <time> |
    Then I should be able to retrieve it by its coordinates:
      | id   | sequence   | time   |
      | <id> | <sequence> | <time> |
    Examples:
      | id                                   | sequence | time                                |
      | 7bf20643-77a7-48b7-b3bd-ac0c2649a8f3 | 0        | 2006-01-02T15:04:05.999999999Z07:00 |

  Scenario Outline: Events with forced sequence numbers should be stored and retrievable
    When I store the event:
      | id   | sequence | time   |
      | <id> | force    | <time> |
    Then I should be able to retrieve it by its coordinates:
      | id   | sequence   | time   |
      | <id> | <sequence> | <time> |
    Examples:
      | id                                   | sequence | time                                |
      | 3d17f9f9-9e76-4fca-8eae-873db34b1c87| 0        | 2006-01-02T15:04:05.999999999Z07:00 |

  Scenario Outline: Event books with mixed forced/fixed should be stored and retrievable
    When I store the event:
      | id   | sequence | time   |
      | <id> | 0        | <time> |
      | <id> | force    | <time> |
      | <id> | 1        | <time> |
    Then I should be able to retrieve it by its coordinates:
      | id   | sequence | time   |
      | <id> | 0        | <time> |
      | <id> | 1        | <time> |
      | <id> | 2        | <time> |
    Examples:
      | id                                   | sequence | time                                |
      | 1c00555a-f335-4757-a842-be3fc3a6b9d5 | 0        | 2006-01-02T15:04:05.999999999Z07:00 |
      | 1c00555a-f335-4757-a842-be3fc3a6b9d5 | 1        | 2006-01-02T15:04:05.999999999Z07:00 |
      | 1c00555a-f335-4757-a842-be3fc3a6b9d5 | 2        | 2006-01-02T15:04:05.999999999Z07:00 |

  Scenario:
    Given a populated database:
      | id                                   | sequence | time                                |
      | 9848ba00-2840-4f6d-ad8a-fdf5cf0cd06c | 0        | 2006-01-02T15:04:05.999999999Z07:00 |
      | 9848ba00-2840-4f6d-ad8a-fdf5cf0cd06c | 1        | 2006-01-02T15:04:05.999999999Z07:00 |
      | 9848ba00-2840-4f6d-ad8a-fdf5cf0cd06c | 2        | 2006-01-02T15:04:05.999999999Z07:00 |
    When I retrieve all events
    Then I should get these events:
      | id                                   | sequence | time                                |
      | 9848ba00-2840-4f6d-ad8a-fdf5cf0cd06c | 0        | 2006-01-02T15:04:05.999999999Z07:00 |
      | 9848ba00-2840-4f6d-ad8a-fdf5cf0cd06c | 1        | 2006-01-02T15:04:05.999999999Z07:00 |
      | 9848ba00-2840-4f6d-ad8a-fdf5cf0cd06c | 2        | 2006-01-02T15:04:05.999999999Z07:00 |

  Scenario:
    Given a populated database:
      | id                                   | sequence | time                                |
      | e84d2250-6a06-4029-8893-037dce50556f | 0        | 2006-01-02T15:04:05.999999999Z07:00 |
      | e84d2250-6a06-4029-8893-037dce50556f | 1        | 2006-01-02T15:04:05.999999999Z07:00 |
      | e84d2250-6a06-4029-8893-037dce50556f | 2        | 2006-01-02T15:04:05.999999999Z07:00 |
    When I retrieve a subset of events starting from value 2
    Then I should get these events:
      | id                                   | sequence | time                                |
      | e84d2250-6a06-4029-8893-037dce50556f | 2        | 2006-01-02T15:04:05.999999999Z07:00 |

  Scenario:
    Given a populated database:
      | id                                   | sequence | time                                |
      | 6e0b80d9-f92c-40a3-aa93-2c9fcfe863ab | 0        | 2006-01-02T15:04:05.999999999Z07:00 |
      | 6e0b80d9-f92c-40a3-aa93-2c9fcfe863ab | 1        | 2006-01-02T15:04:05.999999999Z07:00 |
      | 6e0b80d9-f92c-40a3-aa93-2c9fcfe863ab | 2        | 2006-01-02T15:04:05.999999999Z07:00 |
    When I retrieve a subset of events ending at event 2
    Then I should get these events:
      | id                                   | sequence | time                                |
      | 6e0b80d9-f92c-40a3-aa93-2c9fcfe863ab | 0        | 2006-01-02T15:04:05.999999999Z07:00 |
      | 6e0b80d9-f92c-40a3-aa93-2c9fcfe863ab | 1        | 2006-01-02T15:04:05.999999999Z07:00 |

  Scenario:
    Given a populated database:
      | id                                   | sequence | time                                |
      | 3bdf3e96-87c4-4b8d-89b9-598f29f90133 | 0        | 2006-01-02T15:04:05.999999999Z07:00 |
      | 3bdf3e96-87c4-4b8d-89b9-598f29f90133 | 1        | 2006-01-02T15:04:05.999999999Z07:00 |
      | 3bdf3e96-87c4-4b8d-89b9-598f29f90133 | 2        | 2006-01-02T15:04:05.999999999Z07:00 |
    When I retrieve a subset of events from 1 to 2
    Then I should get these events:
      | id                                   | sequence | time                                |
      | 3bdf3e96-87c4-4b8d-89b9-598f29f90133 | 1        | 2006-01-02T15:04:05.999999999Z07:00 |
