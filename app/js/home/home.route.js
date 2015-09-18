(function() {
    'use strict';

    angular
        .module('app.home')
        .run(appRun);

    appRun.$inject = ['routerHelper'];
    function appRun(routerHelper) {

//         routerHelper.configureStates(getStates(), otherwise);
        routerHelper.configureStates(getStates());
    }

    function getStates() {
        return [
            {
                state: 'home',
                config: {
                    url: '/home',
                    templateUrl: 'js/home/home.html',
                    controller: 'Home',
                    controllerAs: 'vm',
                    title: 'home',
                }
            }
        ];
    }
})();