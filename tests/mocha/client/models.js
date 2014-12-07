if (typeof MochaWeb !== 'undefined'){
    MochaWeb.testOnly(function(){
        describe("Models", function() {
            it("should validate", function() {
                score = new Score();
                chai.assert(check(score, Schemas.Score) === undefined);

                poll = new Poll({"title": 1, "description": "this poll shouldn't pass"});
                try {
                    check(poll, Polls.simpleSchema());
                    chai.assert.fail("incorrect schema passed validation");
                } catch(err) {}

                poll = new Poll({"title": "test", "description": "description here"});
                check(poll, Polls.simpleSchema());
            });
        });
    });
}

