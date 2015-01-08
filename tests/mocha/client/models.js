if (typeof MochaWeb !== 'undefined'){
    MochaWeb.testOnly(function(){
        var userId;

        before(function(done) {
            Accounts.createUser({
                name: "Tester",
                email: "test@example.com",
                password: "testing"
            }, function(err) {
                // TODO: Handle error once user is removed after testing in after()
            });

            Meteor.loginWithPassword("test@example.com", "testing", function(err) {
                if(err) {
                   done(err);
                   return;
                }

                userId = Meteor.userId();
                done();
            });
        });

        after(function(done) {
            // TODO: Make sure user is removed, can't be done client-side.
            // Meteor.users.remove(userId);
            done();
        });

        describe("Models", function() {
            it("should pass validation", function(done) {
                var poll = Poll({
                     title: "test",
                     description: "description here",
                     type: "YesNo"
                });
                chai.assert(Match.test(poll, Polls.simpleSchema()));
                done();
            });

            it("should not pass validation", function(done) {
                var poll = Poll({
                    title: "YOU SHALL NOT PASS!",
                    description: "If Gandalf says so it must be true, but also since the type isn't allowed",
                    type: "InvalidType"
                });
                chai.assert(!Match.test(poll, Polls.simpleSchema()));
                done();
            });
        });
    });
}

