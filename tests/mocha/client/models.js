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
            // TODO: Make sure user is removed
            // Meteor.users.remove(userId);
            done();
        });

        describe("Models", function() {
            it("should validate", function(done) {
                var poll = new Poll({
                     title: "test",
                     description: "description here",
                     type: "YesNo"
                });
                chai.assert.isUndefined(check(poll, Polls.simpleSchema()));

                poll = new Poll({
                    title: 1,
                    description: "this poll shouldn't pass",
                    type: "YesNo"
                });
                try {
                    check(poll, Polls.simpleSchema());
                    chai.assert.fail();
                } catch(e) {
                    if(e instanceof chai.AssertionError) {
                        chai.assert.fail("incorrect schema passed validation");
                    }
                }
                done();
            });
        });
    });
}

