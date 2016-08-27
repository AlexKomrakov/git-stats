angular.module('app', ['ngResource'])
    .config(function ($interpolateProvider) {
        $interpolateProvider.startSymbol('[[').endSymbol(']]')
    })
    .controller('IndexController', function ($resource) {
        this.resource = $resource('/stats');
        this.loadStats = function (query) {
            this.request = this.resource.get(query);
            this.request.$promise.then(function(response){
                this.stats = response;
            }.bind(this));
        }.bind(this);
    });
