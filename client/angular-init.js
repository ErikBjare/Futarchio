app = angular.module('FutarchioApp',['angular-meteor', 'ngResource', 'ngAnimate', 'ngCookies', 'ngSanitize']);

if(Meteor.isClient) {
    Meteor.startup(function() {
        angular.bootstrap(document, ['FutarchioApp']);
    });
}

app.animation('.slide', function() {
    var NG_HIDE_CLASS = 'ng-hide';
    return {
        beforeAddClass: function(element, className, done) {
            if(className === NG_HIDE_CLASS) {
                element.slideUp(done);
            }
        },
        removeClass: function(element, className, done) {
            if(className === NG_HIDE_CLASS) {
                element.hide().slideDown(done);
            }
        }
    };
});
