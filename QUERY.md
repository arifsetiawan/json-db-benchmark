
# MongoDB aggregate

```
db.definitions.aggregate([{
    $match: {
        tenant_id: "google",
        id: "definition:golang:d5ba013a-586e-43c7-9db8-5be52371ee5b"
    }
}, {
    $lookup: {
        "from": "instances",
        "localField": "id",
        "foreignField": "definition",
        "as": "instances"
    }
}, {
    $unwind: {
        path: "$instances",
        preserveNullAndEmptyArrays: true
    }
}, {
    $group: {
        _id: "$_id",
        "data": {
            "$first": "$$ROOT"
        },
        running_count: {
            '$sum': {
                '$cond': [{
                    '$eq': ['$instances.status.running', true]
                }, 1, 0]
            }
        },
        failed_count: {
            '$sum': {
                '$cond': [{
                    '$eq': ['$instances.status.failed', true]
                }, 1, 0]
            }
        },
        completed_count: {
            '$sum': {
                '$cond': [{
                    '$eq': ['$instances.status.completed', true]
                }, 1, 0]
            }
        }
    }
}, {
    $project: {
        id: "$data.id",
        reference: "$data.reference",
        version: "$data.version",
        running_count: "$running_count",
        failed_count: "$failed_count",
        completed_count: "$completed_count"
    }
}]
);
```

# ArangoDB

```
FOR d IN definitions
    FILTER d.tenant_id == "google" && d.id == "definition:golang:1b7678c7-1abd-43ee-bdb2-c063ec668d61"
    LET instances_running = (
        FOR f IN instances
            FILTER f.definition == d.id && f.status.running == true
        RETURN f
    )
    LET instances_completed = (
        FOR f IN instances
            FILTER f.definition == d.id && f.status.completed == true
        RETURN f
    )
    LET instances_failed = (
        FOR f IN instances
            FILTER f.definition == d.id && f.status.failed == true
        RETURN f
    )
    
RETURN {
    id: d.id, 
    reference: d.reference, 
    version: d.version, 
    running_count: LENGTH(instances_running),
    completed_count: LENGTH(instances_completed),
    failed_count: LENGTH(instances_failed)
}
```

# Couchbase

```
SELECT t.id, t.reference, t.version, running_count[0].count AS running_count, failed_count[0].count AS failed_count, completed_count[0].count AS completed_count
FROM engine t 
LET running_count = (
        SELECT count(*) count FROM engine e WHERE REGEXP_LIKE(META(e).id, 'apple:instance:.*') AND e.definition = 'definition:golang:00201be5-ca56-4d7f-abe1-9d526b44e543' AND e.status.running = true
    ),
    failed_count = (
        SELECT count(*) count FROM engine e WHERE REGEXP_LIKE(META(e).id, 'apple:instance:.*') AND e.definition = 'definition:golang:00201be5-ca56-4d7f-abe1-9d526b44e543' AND e.status.failed = true
    ),
    completed_count = (
        SELECT count(*) count FROM engine e WHERE REGEXP_LIKE(META(e).id, 'apple:instance:.*') AND e.definition = 'definition:golang:00201be5-ca56-4d7f-abe1-9d526b44e543' AND e.status.completed = true
    )
    WHERE 
    REGEXP_LIKE(META(t).id, 'apple:definition:.*') AND 
    t.tenant_id IN ['apple'] 
    AND t.id IN ['definition:golang:00201be5-ca56-4d7f-abe1-9d526b44e543']
```

# Postgres

```
SELECT 
    d.definition_id as id, 
    d.data->>'reference' as reference, 
    d.data->>'version' as version, 
    SUM(CASE WHEN i.data @> '{"status":{"running": true}}' THEN 1 ELSE 0 END) AS running_count,
    SUM(CASE WHEN i.data @> '{"status":{"completed": true}}' THEN 1 ELSE 0 END) AS completed_count,
    SUM(CASE WHEN i.data @> '{"status":{"failed": true}}' THEN 1 ELSE 0 END) AS failed_count
FROM engine_definition d 
FULL OUTER JOIN engine_instance i ON i.definition_id = d.definition_id
WHERE
    d.data->>'tenant_id'='google' AND
    d.data->>'id'='definition:golang:587c9db1-db6a-4295-896b-2c6770ea49e5'
GROUP BY d.definition_id 
```

# MySQL

```
SELECT 
    d.definition_id as id, 
    d.data->>"$.reference" as reference, 
    d.data->>"$.version" as version,
    COUNT(CASE WHEN i.data->"$.status.running"=true THEN 1 END) AS running_count,
    COUNT(CASE WHEN i.data->"$.status.completed"=true THEN 1 END) AS completed_count,
    COUNT(CASE WHEN i.data->"$.status.failed"=true THEN 1 END) AS failed_count
FROM engine_definition d 
CROSS JOIN engine_instance i ON i.definition_id = d.definition_id
WHERE
    d.data->>"$.tenant_id"='apple' AND
    d.data->>"$.id"='definition:golang:c0ba0182-d242-4069-b117-d05e4cffa8d8'
GROUP BY d.definition_id 
```