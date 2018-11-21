CREATE TABLE jsonwebtokens 
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
	net     varchar(100) ,
	-- jwt compliance manual link (optional, default stored on IPFS)
	ipfs    varchar(100) ,
	-- predefined CRUD operation permission for specific usecase
	crud    char(60) ,
	-- other required jwt field for usecase
	context text ,
	PRIMARY KEY (jwt_id)
) ENGINE=InnoDB;