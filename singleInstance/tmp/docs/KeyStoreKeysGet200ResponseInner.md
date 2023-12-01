# KeyStoreKeysGet200ResponseInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** |  | 
**Key** | **string** |  | 
**Endpoint** | Pointer to **NullableString** |  | [optional] 
**Provider** | **string** |  | 

## Methods

### NewKeyStoreKeysGet200ResponseInner

`func NewKeyStoreKeysGet200ResponseInner(id int32, key string, provider string, ) *KeyStoreKeysGet200ResponseInner`

NewKeyStoreKeysGet200ResponseInner instantiates a new KeyStoreKeysGet200ResponseInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewKeyStoreKeysGet200ResponseInnerWithDefaults

`func NewKeyStoreKeysGet200ResponseInnerWithDefaults() *KeyStoreKeysGet200ResponseInner`

NewKeyStoreKeysGet200ResponseInnerWithDefaults instantiates a new KeyStoreKeysGet200ResponseInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *KeyStoreKeysGet200ResponseInner) GetId() int32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *KeyStoreKeysGet200ResponseInner) GetIdOk() (*int32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *KeyStoreKeysGet200ResponseInner) SetId(v int32)`

SetId sets Id field to given value.


### GetKey

`func (o *KeyStoreKeysGet200ResponseInner) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *KeyStoreKeysGet200ResponseInner) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *KeyStoreKeysGet200ResponseInner) SetKey(v string)`

SetKey sets Key field to given value.


### GetEndpoint

`func (o *KeyStoreKeysGet200ResponseInner) GetEndpoint() string`

GetEndpoint returns the Endpoint field if non-nil, zero value otherwise.

### GetEndpointOk

`func (o *KeyStoreKeysGet200ResponseInner) GetEndpointOk() (*string, bool)`

GetEndpointOk returns a tuple with the Endpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpoint

`func (o *KeyStoreKeysGet200ResponseInner) SetEndpoint(v string)`

SetEndpoint sets Endpoint field to given value.

### HasEndpoint

`func (o *KeyStoreKeysGet200ResponseInner) HasEndpoint() bool`

HasEndpoint returns a boolean if a field has been set.

### SetEndpointNil

`func (o *KeyStoreKeysGet200ResponseInner) SetEndpointNil(b bool)`

 SetEndpointNil sets the value for Endpoint to be an explicit nil

### UnsetEndpoint
`func (o *KeyStoreKeysGet200ResponseInner) UnsetEndpoint()`

UnsetEndpoint ensures that no value is present for Endpoint, not even an explicit nil
### GetProvider

`func (o *KeyStoreKeysGet200ResponseInner) GetProvider() string`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *KeyStoreKeysGet200ResponseInner) GetProviderOk() (*string, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *KeyStoreKeysGet200ResponseInner) SetProvider(v string)`

SetProvider sets Provider field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


