(function(){
	'use strict';

	angular
		.module('app.core')
		.factory('DataService', DataService);


	DataService.$inject = ['$http', '$location', '$q', exception, logger,'$resource']

	function DataService($http, $location, $q, exception, logger, $resource){

		var service = {
			getLaw: getLaw,

		};

		return service;

		function getLaw(){
			return $resource('/parse');
		}

	}
})