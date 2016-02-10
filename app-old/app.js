var app = angular.module('myapp',["ngResource"]);

app.controller("MainCtl", [ "$scope", "$resource", function($scope, $resource){

	var Law = $resource("/Law/:id", {id: "@id"}, {});

	$scope.books = Law.query();
})