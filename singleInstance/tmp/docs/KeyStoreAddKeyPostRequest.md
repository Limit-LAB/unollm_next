# KeyStoreAddKeyPostRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Key** | **string** |  | 
**Provider** | **string** | ChatGLM | ChatGPT | 
**Endpoint** | Pointer to **NullableString** |  | [optional] 

## Methods

### NewKeyStoreAddKeyPostRequest

`func NewKeyStoreAddKeyPostRequest(key string, provider string, ) *KeyStoreAddKeyPostRequest`

NewKeyStoreAddKeyPostRequest instantiates a new KeyStoreAddKeyPostRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewKeyStoreAddKeyPostRequestWithDefaults

`func NewKeyStoreAddKeyPostRequestWithDefaults() *KeyStoreAddKeyPostRequest`

NewKeyStoreAddKeyPostRequestWithDefaults instantiates a new KeyStoreAddKeyPostRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKey

`func (o *KeyStoreAddKeyPostRequest) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *KeyStoreAddKeyPostRequest) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *KeyStoreAddKeyPostRequest) SetKey(v string)`

SetKey sets Key field to given value.


### GetProvider

`func (o *KeyStoreAddKeyPostRequest) GetProvider() string`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *KeyStoreAddKeyPostRequest) GetProviderOk() (*string, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *KeyStoreAddKeyPostRequest) SetProvider(v string)`

SetProvider sets Provider field to given value.


### GetEndpoint

`func (o *KeyStoreAddKeyPostRequest) GetEndpoint() string`

GetEndpoint returns the Endpoint field if non-nil, zero value otherwise.

### GetEndpointOk

`func (o *KeyStoreAddKeyPostRequest) GetEndpointOk() (*string, bool)`

GetEndpointOk returns a tuple with the Endpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpoint

`func (o *KeyStoreAddKeyPostRequest) SetEndpoint(v string)`

SetEndpoint sets Endpoint field to given value.

### HasEndpoint

`func (o *KeyStoreAddKeyPostRequest) HasEndpoint() bool`

HasEndpoint returns a boolean if a field has been set.

### SetEndpointNil

`func (o *KeyStoreAddKeyPostRequest) SetEndpointNil(b bool)`

 SetEndpointNil sets the value for Endpoint to be an explicit nil

### UnsetEndpoint
`func (o *KeyStoreAddKeyPostRequest) UnsetEndpoint()`

UnsetEndpoint ensures that no value is present for Endpoint, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


