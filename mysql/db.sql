CREATE TABLE credentials 
(
	jwt_id  int 	       NOT NULL AUTO_INCREMENT ,
	-- credentials issuer did
	iss     char(32)       NOT NULL ,
	-- credentials subject for usecase
	sub     char(100)      NOT NULL ,
	-- credentials recipient did
	aud     char(32)       NOT NULL , 
	-- credentials expiration NumericDate(timestamp)
	exp     int UNSIGNED   NOT NULL ,
	-- credentials Not Before NumericDate(timestamp)
	nbf     int UNSIGNED ,
	-- credentials Issued At NumericDate(timestamp)
	iat     int UNSIGNED ,
	-- jwt identifier to avoid Replay attack for specific usecase
	jti  	varchar(255) ,
	-- blockchain identifier, defalut is eth_ropsten
	net     varchar(100)   NOT NULL ,
	-- jwt compliance manual link (optional, default stored on IPFS)
	ipfs    varchar(100) ,
	-- predefined CRUD operation permission for specific usecase
	-- crud    char(60)       NOT NULL ,
	-- other required jwt field for usecase
	context text ,
	-- base64-url encoded jwt
	credential text ,
	-- jwt permission status
	status  int            NOT NULL ,
	PRIMARY KEY (jwt_id)
) ENGINE=InnoDB;

CREATE TABLE updated_credentials 
(
	jwt_id  int 	       NOT NULL AUTO_INCREMENT ,
	-- credentials issuer did
	iss     char(32)       NOT NULL ,
	-- credentials subject for usecase
	sub     char(100)      NOT NULL ,
	-- credentials recipient did
	aud     char(32)       NOT NULL , 
	-- credentials expiration NumericDate(timestamp)
	exp     int UNSIGNED   NOT NULL ,
	-- credentials Not Before NumericDate(timestamp)
	nbf     int UNSIGNED ,
	-- credentials Issued At NumericDate(timestamp)
	iat     int UNSIGNED ,
	-- jwt identifier to avoid Replay attack for specific usecase
	jti  	varchar(255) ,
	-- blockchain identifier, defalut is eth_ropsten
	net     varchar(100)   NOT NULL ,
	-- jwt compliance manual link (optional, default stored on IPFS)
	ipfs    varchar(100) ,
	-- predefined CRUD operation permission for specific usecase
	-- crud    char(60)       NOT NULL ,
	-- other required jwt field for usecase
	context text ,
	-- base64-url encoded jwt
	credential text ,
	-- jwt permission status
	status  int            NOT NULL ,
	PRIMARY KEY (jwt_id)
) ENGINE=InnoDB;