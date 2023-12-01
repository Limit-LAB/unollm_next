# KeyStoreMapToPostRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Keys** | **[]int32** | APIKey的ID | 
**MapTo** | **int32** | 目标Key | 

## Methods

### NewKeyStoreMapToPostRequest

`func NewKeyStoreMapToPostRequest(keys []int32, mapTo int32, ) *KeyStoreMapToPostRequest`

NewKeyStoreMapToPostRequest instantiates a new KeyStoreMapToPostRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewKeyStoreMapToPostRequestWithDefaults

`func NewKeyStoreMapToPostRequestWithDefaults() *KeyStoreMapToPostRequest`

NewKeyStoreMapToPostRequestWithDefaults instantiates a new KeyStoreMapToPostRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeys

`func (o *KeyStoreMapToPostRequest) GetKeys() []int32`

GetKeys returns the Keys field if non-nil, zero value otherwise.

### GetKeysOk

`func (o *KeyStoreMapToPostRequest) GetKeysOk() (*[]int32, bool)`

GetKeysOk returns a tuple with the Keys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeys

`func (o *KeyStoreMapToPostRequest) SetKeys(v []int32)`

SetKeys sets Keys field to given value.


### GetMapTo

`func (o *KeyStoreMapToPostRequest) GetMapTo() int32`

GetMapTo returns the MapTo field if non-nil, zero value otherwise.

### GetMapToOk

`func (o *KeyStoreMapToPostRequest) GetMapToOk() (*int32, bool)`

GetMapToOk returns a tuple with the MapTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMapTo

`func (o *KeyStoreMapToPostRequest) SetMapTo(v int32)`

SetMapTo sets MapTo field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


