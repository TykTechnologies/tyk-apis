package test

//+o:title=Tyk GW
//+o:version=v3.0.4
//+o:license:name="MIT",url=""

//+o:security={Oauth2:scope2;scope1}
//+o:security={openid:scopea;scopeb}

//+o:tags:name=pets,description="Everything about your Pets",externalDocs={url:"http://docs.my-api.com/pet-operations.htm"}
//+o:tags:name=store,description="Access to Petstore orders",externalDocs={url:"http://docs.my-api.com/store-orders.htm"}
