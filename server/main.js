Router.route('/api/0/test', function () {
    var req = this.request;
    var res = this.response;
    res.end(JSON.stringify({msg: 'hello from the server'}));
}, {where: 'server'});

Router.route('/api/0/polls', function () {
    var req = this.request;
    var res = this.response;
    console.log(Polls); 
    res.end(JSON.stringify([1,2,3,"test"]));
}, {where: 'server'});

Meteor.startup(function () {
    // code to run on server at startup
});
