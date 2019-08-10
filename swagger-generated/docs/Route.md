# Route

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Athlete** | [***SummaryAthlete**](SummaryAthlete.md) |  | [optional] [default to null]
**Description** | **string** | The description of the route | [optional] [default to null]
**Distance** | **float32** | The route&#39;s distance, in meters | [optional] [default to null]
**ElevationGain** | **float32** | The route&#39;s elevation gain. | [optional] [default to null]
**Id** | **int32** | The unique identifier of this route | [optional] [default to null]
**Map_** | [***PolylineMap**](PolylineMap.md) |  | [optional] [default to null]
**Name** | **string** | The name of this route | [optional] [default to null]
**Private** | **bool** | Whether this route is private | [optional] [default to null]
**Starred** | **bool** | Whether this route is starred by the logged-in athlete | [optional] [default to null]
**Timestamp** | **int32** |  | [optional] [default to null]
**Type_** | **int32** | This route&#39;s type (1 for ride, 2 for runs) | [optional] [default to null]
**SubType** | **int32** | This route&#39;s sub-type (1 for road, 2 for mountain bike, 3 for cross, 4 for trail, 5 for mixed) | [optional] [default to null]
**Segments** | [**[]SummarySegment**](SummarySegment.md) | The segments traversed by this route | [optional] [default to null]
**Directions** | [**[]RouteDirection**](RouteDirection.md) | The directions of this route | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


