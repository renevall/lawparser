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
                    templateUrl: 'js/law/review.html',
                    controller: 'ReviewLaw',
                    controllerAs: 'vm',
                    title: 'review law',
                    settings: {
                        nav: 1,
                        content: '<i class="fa fa-dashboard"></i> Review Law'
                    }
                }
            },
            {
                state: 'upload',
                config:{
                    url: '/upload',
                    templateUrl: 'js/law/upload.html',
                    controller: 'UploadController',
                    controllerAs: 'vm',
                    title: 'Upload Law'
                }
            }
        ];
    }
})();