/**
 * The file provides stubs for JavaScript objects accessible from HTTP Client response handler scripts.
 * It doesn't perform any real operation and should be used for documentation purpose.
 */

/**
 * The object holds HTTP Client session meta data, e.g. list of global variables.
 *
 * @type {HttpClient}
 */
var client = new HttpClient();

/**
 * The object holds information about HTTP Response.
 *
 * @type {HttpResponse}
 */
var response = new HttpResponse();

/**
 * HTTP Client session meta data, e.g. list of global variables.
 *
 * HTTP Client session is started on IDE start and ends on IDE close,
 * values are not preserved after IDE restart.
 */
function HttpClient() {
  /**
   * Global variables defined in response handler scripts,
   * can be used as variables in HTTP Requests,
   *
   * Example:
   * ### Authorization request, receives token as an attribute of json body
   * GET https://example.com/auth
   *
   * > {% client.global.set("auth_token", response.body.token) %}
   *
   * ### Request executed with received auth_token
   * GET http://example.com/get
   * Authorization: Bearer {{auth_token}}
   *
   * @type {Variables}
   */
  this.global = new Variables();

  /**
   * Creates test with name 'testName' and body 'func'.
   * All tests will be executed right after response handler script.
   * @param testName {string}
   * @param func {function}
   */
  this.test = function (testName, func) {
  };

  /**
   * Checks that condition is true and throw an exception otherwise.
   * @param condition {boolean}
   * @param message {string} optional parameter, if specified it will be used as an exception message.
   */
  this.assert = function (condition, message) {
  };

  /**
   * Prints text to the response handler or test stdout and then terminates the line.
   */
  this.log = function (text) {
  };
}

/**
 * Variables storage, can be used to define, undefine or retrieve variables.
 */
function Variables() {
  /**
   * Saves variable with name 'varName' and sets its value to 'varValue'.
   * @param varName {string}
   * @param varValue {string}
   */
  this.set = function (varName, varValue) {
  };

  /**
   * Returns value of variable 'varName'.
   * @param varName {string}
   * @returns {string}
   */
  this.get = function (varName) {
    return varValue
  };

  /**
   * Checks no variables are defined.
   * @returns {boolean}
   */
  this.isEmpty = function () {
    return true
  };

  /**
   * Removes variable 'varName'.
   * @param varName {string}
   */
  this.clear = function (varName) {
  };

  /**
   * Removes all variables.
   */
  this.clearAll = function () {
  };
}

/**
 * HTTP Response data object, contains information about response content, headers, status, etc.
 */
function HttpResponse() {
  /**
   * Response content, it is a string or JSON object if response content-type is json.
   * @type {string|object}
   */
  this.body = " ";

  /**
   * Response headers storage.
   * @type {ResponseHeaders}
   */
  this.headers = new ResponseHeaders();

  /**
   * Response status, e.g. 200, 404, etc.
   * @type {int}
   */
  this.status = 200;

  /**
   * Value of 'Content-Type' response header.
   * @type {ContentType}
   */
  this.contentType = new ContentType
}

/**
 * Headers storage, can be use to retrieve data about header value.
 */
function ResponseHeaders() {
  /**
   * Retrieves the first value of 'headerName' response header or null otherwise.
   * @param headerName {string}
   * @returns {string|null}
   */
  this.valueOf = function (headerName) {
    return headerValue
  };

  /**
   * Retrieves all values of 'headerName' response header. Returns empty list if header with 'headerName' doesn't exist.
   * @param headerName {string}
   * @returns {string[]}
   */
  this.valuesOf = function (headerName) {
    return headerValue
  };
}

/**
 * Content type data object, contains information from 'Content-Type' response header.
 */
function ContentType() {
  /**
   * MIME type of the response,
   * e.g. 'text/plain', 'text/xml', 'application/json'.
   * @type {string}
   */
  this.mimeType = "application/json";

  /**
   * String representation of the response charset,
   * e.g. utf-8.
   * @type {string}
   */
  this.charset = "utf-8";
}
