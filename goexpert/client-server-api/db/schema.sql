USE goexpert;
CREATE TABLE if not exists cotations (
    id integer not null auto_increment,
    code varchar(3),
    codein varchar(3),
    name varchar(50),
    high varchar(7),
    low varchar(7),
    varBid varchar(7),
    pctChange varchar(4),
    bid varchar(7),
    ask varchar(7),
    timestamp varchar(15),
    createDate varchar(20),
    PRIMARY KEY (id)
);
