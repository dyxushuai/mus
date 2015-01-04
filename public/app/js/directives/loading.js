/**
 * Loading Directive
 * @see http://tobiasahlin.com/spinkit/
 */

module.exports = function () {
    var directive = {
        restrict: 'AE',
        template: '<div class="loading"><div class="dot1"></div><div class="dot2"></div></div>'
    };
    return directive;
};