#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

@all
@vault_server
Feature: Vault server API

  @vault_server_healthcheck
  Scenario: Vault server healthcheck
    When an HTTP GET is sent to "https://localhost:9099/healthcheck"
    Then the JSON path "status" of the response equals "success"

  @vault_server_create
  Scenario: Creates a vault
    When Create a new vault using the vault server "https://localhost:9099"
    Then Save a document with the following id "M3aS9xwj8ybCwHkEiCJJR1"
    And Check that a document with id "M3aS9xwj8ybCwHkEiCJJR1" is stored

  @vault_server_create_no_id
  Scenario: Creates a vault without id
    When Create a new vault using the vault server "https://localhost:9099"
    Then Save a document without id and save the result ID as "created"
    And Check that a document with id "created" is stored

  @vault_server_create_authorization
  Scenario: Creates an authorization
    When Create a new vault using the vault server "https://localhost:9099"
    Then Save a document with the following id "M3aS9xwj8ybCwHkEiCJJR1"
    And Check that a document with id "M3aS9xwj8ybCwHkEiCJJR1" is stored
    Then Create a new authorization and save the result as "auth"
    And Check that a document with id "M3aS9xwj8ybCwHkEiCJJR1" is available for "auth"