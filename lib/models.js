Schemas = {};

Schemas.Score = new SimpleSchema({
    up: {
        type: Number,
        min: 0
    },
    down: {
        type: Number,
        min: 0
    }
});

Schemas.Poll = new SimpleSchema({
    title: {
        type: String
    },
    description: {
        type: String
    },
    score: {
        type: Schemas.Score
    }
});

Schemas.Statement = new SimpleSchema({
    title: {
        type: String
    },
    description: {
        type: String
    },
    score: {
        type: Schemas.Score
    }
});

Score = function() {
    return {"up": 0, "down": 0};
};

Polls = new Mongo.Collection("polls");
Polls.attachSchema(Schemas.Poll);
Poll = function(poll) {
    poll.score = new Score();
    return poll;
};

Statements = new Mongo.Collection("statements");
Statements.attachSchema(Schemas.Statement);
Statement = function(stmt) {
    stmt.score = new Score();
    return stmt;
};

