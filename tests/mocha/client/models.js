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
                var stmt = Statement({
                     title: "test",
                     description: "description here"
                });
                chai.assert(Match.test(stmt, Statements.simpleSchema()));
                done();

            });

            it("should not pass validation", function(done) {
                var stmt = Statement({
                    description: "This won't pass validation because it misses a title"
                });
                chai.assert(!Match.test(stmt, Statements.simpleSchema()));

                stmt = Statement({
                    title: "This won't pass validation because it misses a description"
                });
                chai.assert(!Match.test(stmt, Statements.simpleSchema()));

                done();
            });
        });
    });
}

