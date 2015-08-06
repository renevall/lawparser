( function(){

	'use strict';
	angular
		.module('app.law')
		.run(appRun);

	function appRun(routerhelper){
		routerhelper.configureRoutes(getRoutes());

	}

	function getRoutes(){
		return [
			{
				url: '/review',
				config: {
					templateUrl: 'app/law/rewiew.html',
					controller: 'Review',
					controllerAs: 'vm',
					title: 'review law',
					settings: {
						nav: 2,
						content: '<i class="fa fa-lock"></i> Review Law'
					}
				}
			}
		];
	}
})