(function() {
    'use strict';

    angular
        .module('app.law')
        .run(appRun);

    appRun.$inject = ['routerHelper'];
    function appRun(routerHelper) {
        console.log();
        routerHelper.configureStates(getStates());
    }

    function getStates() {
        return [
            {
                state: 'review',
                config: {
                    url: '/review',
                    templateUrl: 'review.html',
                    controller: 'ReviewLaw',
                    controllerAs: 'vm',
                    title: 'review law',
                    settings: {
                        nav: 1,
                        content: '<i class="fa fa-dashboard"></i> Review Law'
                    }
                }
            }
        ];
    }
})();