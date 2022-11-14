Feature: API
  Implements pinning API

  Scenario: No pins
    Given a pinning server
    And an empty database
    When a list of pins is requested
    Then the response code should be 200
    And there should be 0 pins in result

  Scenario: Create a pin
    Given a pinning server
    When a pin is created with name "test" and CID "Qma71JMRwZc2aVMZ5McmbggfTgMJJQ8k3HKM8GpMeBR2CU"
    Then the response code should be 202
    And a list of pins is requested
    And the response code should be 200
    And there should be 1 pins in result

  Scenario: Create a pin with a bad CID
    Given a pinning server
    When a pin is created with name "baaaad" and CID "Qma"
    Then the response code should be 400
    And a list of pins is requested
    And the response code should be 200
    And there should be 1 pins in result

  Scenario: Delete pin
    Given a pinning server
    And an empty database
    And a pin is created with name "deleteme" and CID "Qma71JMRwZc2aVMZ5McmbggfTgMJJQ8k3HKM8GpMeBR2CU"
    When the pin is deleted
    And the response code should be 202
    And a list of pins is requested
    And the response code should be 200
    And there should be 0 pins in result