(function() {
    'use strict';

    var core = angular.module('app.core');

    core.config(toastrConfig);

    /* @ngInject */
    toastrConfig.$inject = ['toastr'];
    function toastrConfig(toastr) {
        toastr.options.timeOut = 4000;
        toastr.options.positionClass = 'toast-bottom-right';
    }

    var config = {
        appErrorPrefix: '[GulpPatterns Error] ', //Configure the exceptionHandler decorator
        appTitle: 'Gulp Patterns Demo',
        imageBasePath: '/content/images/photos/',
        unknownPersonImageSource: 'unknown_person.jpg',
        version: '1.0.0'
    };

    core.value('config', config);

    core.config(configure);

    /* @ngInject */
    configure.$inject =['$logProvider', 'routerHelperProvider', 'exceptionHandlerProvider']
    function configure ($logProvider, routerHelperProvider, exceptionHandlerProvider) {
        // turn debugging off/on (no info or warn)
        if ($logProvider.debugEnabled) {
            $logProvider.debugEnabled(true);
        }
        exceptionHandlerProvider.configure(config.appErrorPrefix);
        configureStateHelper();

        ////////////////

        function configureStateHelper() {
            var resolveAlways = { /* @ngInject */
                ready: ['dataservice',function(dataservice) {
                    //return dataservice.ready();

                    return dataservice.ready();
                }]
            };

            routerHelperProvider.configure({
                docTitle: 'NG-Modular: ',
                resolveAlways: resolveAlways
            });
        }
    }
})();