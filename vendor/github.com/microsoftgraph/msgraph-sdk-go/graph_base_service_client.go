package msgraphsdkgo

import (
    i25911dc319edd61cbac496af7eab5ef20b6069a42515e22ec6a9bc97bf598488 "github.com/microsoft/kiota-serialization-json-go"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i7294a22093d408fdca300f11b81a887d89c47b764af06c8b803e2323973fdb83 "github.com/microsoft/kiota-serialization-text-go"
    i009f47bbce65ccdb7303730eed71e6bab3ae2f8e4e918bc9e94341d28624af97 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
    i07d47a144340607d6d6dbd93575e531530e4f1cc6091c947ea0766f7951ffd34 "github.com/microsoftgraph/msgraph-sdk-go/shares"
    i088cccaaaff488138f258ec600ca804c5dc9548772ebe52bf3cb7a3eaf4b9fdf "github.com/microsoftgraph/msgraph-sdk-go/workbooks"
    i0906e75d8a44bf92212e084e1d2f62d03887dcec6a5c8535e92ccc04c1e5fdec "github.com/microsoftgraph/msgraph-sdk-go/solutions"
    i185698f71f6301975f0627ee999e6e91920d8fa9c00bdef3487b9f349e2df04e "github.com/microsoftgraph/msgraph-sdk-go/directoryobjects"
    i1a1369b1521a8ac4885166fd68eae4247248a891006fea464d2eea2a271b2cdb "github.com/microsoftgraph/msgraph-sdk-go/permissiongrants"
    i1b75be7b5675627960b4672ab148be21ff379d5cbc0e62f6bc5b97d54464f8b5 "github.com/microsoftgraph/msgraph-sdk-go/teamstemplates"
    i1be0f1b1da466bc62355d411ef490acbd8dc0ec5ca4d3448c7eb73e5caffafc3 "github.com/microsoftgraph/msgraph-sdk-go/education"
    i1d6652ecc686b20c37a9a3448b26db8187e284e1a4017cab8876b02b97557436 "github.com/microsoftgraph/msgraph-sdk-go/grouplifecyclepolicies"
    i1dc06c4b7f499cb445a6c55e466abd6d7466bb35a2683c675909db23c57898e7 "github.com/microsoftgraph/msgraph-sdk-go/authenticationmethodconfigurations"
    i20b08d3949f1191430a14a315e0758a1f131dc59bbdc93e654f1dd447a6af14c "github.com/microsoftgraph/msgraph-sdk-go/auditlogs"
    i286f3babd79fe9ec3b0f52b6ed5910842c0adaeff02be1843d0e01c56d9ba6d9 "github.com/microsoftgraph/msgraph-sdk-go/search"
    i2a252d42835bdab6d88bf938595da6cf029001f9ca970d6f599cecf0ca27f8e5 "github.com/microsoftgraph/msgraph-sdk-go/directoryroletemplates"
    i32d45c1243c349600fbe53b2f9641bb59857a3326037587cbe4e347b46ad207e "github.com/microsoftgraph/msgraph-sdk-go/identitygovernance"
    i35d7bbcc8f7e8b8e9525ea0ee5b3c51c3a1a58f9ed512b727d181bfcd08eb032 "github.com/microsoftgraph/msgraph-sdk-go/security"
    i3e9b5129e2bb8b32b0374f7afe2536be6674d73df6c41d7c529f5a5432c4e0aa "github.com/microsoftgraph/msgraph-sdk-go/agreementacceptances"
    i4794c103c0d044c27a3ca3af0a0e498e93a9863420c1a4e7a29ef37590053c7b "github.com/microsoftgraph/msgraph-sdk-go/groupsettings"
    i4a624e38d68c2a9fc4db1ea915bcaffde116f967f58ec2c99e2ea8bbff3690e1 "github.com/microsoftgraph/msgraph-sdk-go/schemaextensions"
    i4ac7f0a844871066493521918f268cafe2a25c71c28a98221ea3f22d5153090f "github.com/microsoftgraph/msgraph-sdk-go/policies"
    i4c91eeb51f03f9d59a342065f7c6ee027ad1fe84ada6b1946b8162c5ae146cfb "github.com/microsoftgraph/msgraph-sdk-go/devices"
    i51b9802eedc1a25686534d117657be902df58c07e90ac6ea84501100998084d9 "github.com/microsoftgraph/msgraph-sdk-go/communications"
    i5310ba7d4cfddbf5de4c1be94a30f9ca8c747c30a87e76587ce88d1cbfff01b4 "github.com/microsoftgraph/msgraph-sdk-go/applicationtemplates"
    i535d6c02ba98f73ff3a8c1c12a035ba5de51606f93aa2c0babdfed56fe505550 "github.com/microsoftgraph/msgraph-sdk-go/certificatebasedauthconfiguration"
    i58857a108d6e260e56ef0dd7e783668388f113eb436006780703ac59f0abb3b1 "github.com/microsoftgraph/msgraph-sdk-go/privacy"
    i61686672307beee899fe5a14188df42982da47730f55a14800b102cd10ab2d72 "github.com/microsoftgraph/msgraph-sdk-go/localizations"
    i62c2771f3f3a1e5e085aedcde54473e9f043cc57b9ce4dd88980a77aca7a5a10 "github.com/microsoftgraph/msgraph-sdk-go/identityproviders"
    i638650494f9db477daff56d31ff923f5c100f72df0257ed7fa5c222cb1a77a94 "github.com/microsoftgraph/msgraph-sdk-go/deviceappmanagement"
    i663c30678b300c2c4b619c4964b4326e471e4da61a44d7c39f752349da7a468e "github.com/microsoftgraph/msgraph-sdk-go/identityprotection"
    i6bf2d83eea06710580ad0d54b886ac4e14cbab0d1d84937f340f02b99f8f5738 "github.com/microsoftgraph/msgraph-sdk-go/reports"
    i71117da372286e863c042a526ec1361696ab14b838a5b77db5bc54386d436543 "github.com/microsoftgraph/msgraph-sdk-go/me"
    i738daeb889f22c1e163aee5a37a094b55b1d815dc76d4802d64e4e1b2e44206c "github.com/microsoftgraph/msgraph-sdk-go/devicemanagement"
    i79ca23a9ac0659e1330dd29e049fe157787d5af6695ead2ff8263396db68d027 "github.com/microsoftgraph/msgraph-sdk-go/identity"
    i7c9d1b36ac198368c1d8bed014b43e2a518b170ee45bf02c8bbe64544a50539a "github.com/microsoftgraph/msgraph-sdk-go/admin"
    i7d140130aac6882792a019b5ebe51fe8d69dfd63ec213c2e3cd98282ce2d0428 "github.com/microsoftgraph/msgraph-sdk-go/appcatalogs"
    i86cada4d4a5f2f8a9d1e7a85eacd70a661ea7b20d2737008c0719e95b5be3e16 "github.com/microsoftgraph/msgraph-sdk-go/oauth2permissiongrants"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    i926bd489c52af20f44aacc8a450bb0a062290f1d1e44c2fe78d6cc1595c12524 "github.com/microsoftgraph/msgraph-sdk-go/drive"
    i93194122344a685a2f9264205dc6d89a5ba39afdcea57fd0ade8f54b6f137c02 "github.com/microsoftgraph/msgraph-sdk-go/applications"
    i9429d7aae2f5c1dabbecc9411e8ad2b733d29338bc0c0436eeccc94605c461b7 "github.com/microsoftgraph/msgraph-sdk-go/print"
    i957076b10ba162b23efec7b94dd26b84c6475d285449c1cbc9c5b85910d36a12 "github.com/microsoftgraph/msgraph-sdk-go/domains"
    ia3e0f7c2d21d5c73ecb8a7552177d0fe444ae0522290dd1c4b5559e449b118af "github.com/microsoftgraph/msgraph-sdk-go/places"
    ia4b736f581ceef30e9ef8cebd9a6c2b932628e087982ff3dd2c9a0f1a920a918 "github.com/microsoftgraph/msgraph-sdk-go/compliance"
    ia6e876e3ed2d92c29c13dbc8c37513bc38d0d5f05ab9321e43a25ff336912a2d "github.com/microsoftgraph/msgraph-sdk-go/groups"
    iaca6694a878291d0e4021155b406c19d3080cdfc382b456e43c71264d4d9e519 "github.com/microsoftgraph/msgraph-sdk-go/domaindnsrecords"
    ib14d748b564c787931c10f1c7ba6856eeddea29a5b9e5c5c27eb1224ff65e5c4 "github.com/microsoftgraph/msgraph-sdk-go/directory"
    ib3217193884e00033cb8182cac52178dfa3b20ce9c4eb48e37a6217882d956ae "github.com/microsoftgraph/msgraph-sdk-go/external"
    ib33fc5e9889e020c0c572578957f59819123a589c61fd7f3eb37eb7958b525ee "github.com/microsoftgraph/msgraph-sdk-go/datapolicyoperations"
    ib68fa8e66bda853b3a33c491e8a66ca665897dab129192b2c97289266c4a1415 "github.com/microsoftgraph/msgraph-sdk-go/informationprotection"
    ibaef614e7692eebc6aaa8080b8ac29169fdf539f24925bc1de4465a3fcdac177 "github.com/microsoftgraph/msgraph-sdk-go/chats"
    ic5e701d75e87f15ce153687b00984a314f7eeea8cfdc77cd9ad648e5ccbc7fbd "github.com/microsoftgraph/msgraph-sdk-go/invitations"
    ic949a0bb5066d68760e8502a7f9db83f571d9e01e38fad4aadf7268188e52df0 "github.com/microsoftgraph/msgraph-sdk-go/organization"
    icabdee72951e77325f237b36d388a199c87e65f67652b6bb85723aba847d7e83 "github.com/microsoftgraph/msgraph-sdk-go/connections"
    ice10f31b9db59ba91184d2b882172edb754f885050cf0830aa2b7c8ff880556b "github.com/microsoftgraph/msgraph-sdk-go/scopedrolememberships"
    id007bc768abbff1131aab64890cdcd0411159a946e9df27140c5f7cf8f249647 "github.com/microsoftgraph/msgraph-sdk-go/subscribedskus"
    id2ac823944414906187dbe4e6ca3b5e46886b9db738d2c1c27de6df8b1bebd61 "github.com/microsoftgraph/msgraph-sdk-go/groupsettingtemplates"
    id4615a956cb1e7edabf8f5a4bc131d1ceca9a13d0f79ae0e122997452a9a0a4e "github.com/microsoftgraph/msgraph-sdk-go/directoryroles"
    id81f15a01b3ceaefa8b1b55f4ee944912f2179aafc4d873f0a2eaf0853eeccd0 "github.com/microsoftgraph/msgraph-sdk-go/authenticationmethodspolicy"
    idb8230b65f4a369c23b4d9b41ebe568c657c92f8f77fe36d16d64528b3a317a3 "github.com/microsoftgraph/msgraph-sdk-go/subscriptions"
    ie05ac24b652f7d895cca374316c093c4ca40dd2df0f1518c465233d6432b1ef9 "github.com/microsoftgraph/msgraph-sdk-go/teamwork"
    ie3631868038c44f490dbc03525ac7249d0523c29cc45cbb25b2aebcf470d6c0c "github.com/microsoftgraph/msgraph-sdk-go/contracts"
    ie66b913c1bc1c536bc8db5d185910e9318f621374e016f95e36e9d59b7127f63 "github.com/microsoftgraph/msgraph-sdk-go/planner"
    ieaa2790c8b7fa361674e69e4a385e279c8c641adf79d86e5b0ca566591a507e8 "github.com/microsoftgraph/msgraph-sdk-go/agreements"
    iefc72d8a17962d4db125c50866617eaa15d662c6e3fb13735d477380dcc0dbe3 "github.com/microsoftgraph/msgraph-sdk-go/drives"
    if398f5c2f1cb53106e045240edd469d82f1854899fd95cfdf8f559b19375750c "github.com/microsoftgraph/msgraph-sdk-go/branding"
    if39bc788926a05e976b265ecfc616408ca12af399df9ce3a2bb348fe89708057 "github.com/microsoftgraph/msgraph-sdk-go/teams"
    if51cca2652371587dbc02e65260e291435a6a8f7f2ffb419f26c3b9d2a033f57 "github.com/microsoftgraph/msgraph-sdk-go/contacts"
    if5372351befdb652f617b1ee71fbf092fa8dd2a161ba9c021bc265628b6ea82b "github.com/microsoftgraph/msgraph-sdk-go/sites"
    if5555fa41b6637688bcf8c25c62a041258f4dc6eacb38ad42d91c66f222ee182 "github.com/microsoftgraph/msgraph-sdk-go/rolemanagement"
    if6ffd1464db2d9c22e351b03e4c00ebd24a5353cd70ffb7f56cfad1c3ceec329 "github.com/microsoftgraph/msgraph-sdk-go/users"
    i009581390843c78f63b06f9dcefeeb5ef2a124a2ac1dcbad3adbe4d0d5e650af "github.com/microsoftgraph/msgraph-sdk-go/users/item"
    i0832d6fdb5fdd25a7a1f3a0a99c90d6ef5c44fdc1f1e1f72f7d777d32cf5ea19 "github.com/microsoftgraph/msgraph-sdk-go/agreements/item"
    i120b7d5508b5c9e9e49c562cc2c54282d0cac65c8ec72e8928f45cc697956915 "github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals/item"
    i16322b5a209cc0bd8b212f5c6ff6e00e3c536005637b4e1beb5a7076653b687a "github.com/microsoftgraph/msgraph-sdk-go/agreementacceptances/item"
    i164d67321f322030ad927126612d90d5880d461e3357bd32611832679c00aff2 "github.com/microsoftgraph/msgraph-sdk-go/domains/item"
    i1a7ecef3f68d1e375afc82029319b84e73d0840cf0fa231262dfbe877c5b0f30 "github.com/microsoftgraph/msgraph-sdk-go/shares/item"
    i1cb1089a14118ff1e1e1bd2d5f465b91b163bfc288e9bb57dc21502014e6b0c1 "github.com/microsoftgraph/msgraph-sdk-go/applications/item"
    i206b88e4abcec8b25d993b477b59473f3e9420358bcd878d07c5ed2b531ccf4c "github.com/microsoftgraph/msgraph-sdk-go/datapolicyoperations/item"
    i23bab38fb8688d4bab0b6ffc533eb085d40e58af49a27ab228a8d1ad3e0ab203 "github.com/microsoftgraph/msgraph-sdk-go/permissiongrants/item"
    i24a463345a902b042d0fa0b40e03cab9b230ab6328a113241a23ab1b81c0bcd1 "github.com/microsoftgraph/msgraph-sdk-go/localizations/item"
    i292f565d97fd74d4ba58ae7a10fd84e8982b27c0e3ba8747d8fd7f8b61482271 "github.com/microsoftgraph/msgraph-sdk-go/identityproviders/item"
    i2a87b2a2ef0c6367c3abc2e909d9176a42e3124d1896753ad0a25d9a6a881c32 "github.com/microsoftgraph/msgraph-sdk-go/teamstemplates/item"
    i2cc346a33133d41934cfa6e862c4eb0d4cf1dc36485198c479852b282338a897 "github.com/microsoftgraph/msgraph-sdk-go/applicationtemplates/item"
    i31a01e7736f53f5c2e0acd0c456ed6dcc6065ca402aaf31b6a1526bd3bf57394 "github.com/microsoftgraph/msgraph-sdk-go/workbooks/item"
    i31ea098bea3bcc5d0f40d4471abad6488910c0e4872783682f25fb2bdbcffdb6 "github.com/microsoftgraph/msgraph-sdk-go/schemaextensions/item"
    i3ea73b2c3d03b959a8d3906c2b1b073951feb694dd6984ab1eea4e6b8c1d0858 "github.com/microsoftgraph/msgraph-sdk-go/contacts/item"
    i4b66adb64e72d30468c3d30826f1df6386f1703ce53cf4a46d2a6cbbfd88b747 "github.com/microsoftgraph/msgraph-sdk-go/teams/item"
    i4c038322c9990a1737934ea712b7d72201832cc50fbe4b6e382968270c3611c3 "github.com/microsoftgraph/msgraph-sdk-go/chats/item"
    i55644ead1eb4b291341a5935abcd8425f7456cccabc4594ef4aee967d51fc863 "github.com/microsoftgraph/msgraph-sdk-go/domaindnsrecords/item"
    i62e63d0f24e0caea5c3b2202a0ee0923bdf496f3e118faa4ee49895e02eecfde "github.com/microsoftgraph/msgraph-sdk-go/invitations/item"
    i728b903f252feeed28263ff4e0da95a1d0f7c831116db07abb00de41db959892 "github.com/microsoftgraph/msgraph-sdk-go/connections/item"
    i7bd49b9f046ac5c4108447b2999dd223be9f057c9910bdc197bcba51be5ac641 "github.com/microsoftgraph/msgraph-sdk-go/sites/item"
    i838f22f4ef3405018163f34ed56f20e99e7bb1a5bd0bfd7009a72d054a09a36c "github.com/microsoftgraph/msgraph-sdk-go/subscribedskus/item"
    i861ab10f93223993f014e54921d0feda5bc5dc9a8996dbbbb728e672d3e8162e "github.com/microsoftgraph/msgraph-sdk-go/certificatebasedauthconfiguration/item"
    i86bb3d05e1a6bbdb496bd3c65829f1a6eb272be42e9ac6060f873dfbd921e4ea "github.com/microsoftgraph/msgraph-sdk-go/contracts/item"
    i960f074bae2d1f849aec23c162b7be41055a1485f8efd075e3c89717cc4ac8f5 "github.com/microsoftgraph/msgraph-sdk-go/directoryroles/item"
    ia5422e2deae41871358311d10a7b0d4a60e914828d2fe80bf0a1bd96c1ff2a2f "github.com/microsoftgraph/msgraph-sdk-go/organization/item"
    ib6d66da0f7d4860b7205f5fdb1200fc9000adb4fbc853a2f05f70c644580220f "github.com/microsoftgraph/msgraph-sdk-go/devices/item"
    ic73c2557206a32cd6d6e58acca928163e176fe1cf7382bdae4f55d28ff56e345 "github.com/microsoftgraph/msgraph-sdk-go/drives/item"
    icd43f65a6cebd6c2685c244d7f46f49a951d8647717cded612ba79705e6aa7c7 "github.com/microsoftgraph/msgraph-sdk-go/grouplifecyclepolicies/item"
    icf62d3bb4e29c8437041430705851ef556cb3cf91d77df26e8eaf92a32e9814e "github.com/microsoftgraph/msgraph-sdk-go/directoryroletemplates/item"
    id5e9a05bae8f5cd30c57fd87690f009f004424eafeb45208f44e7ed46f8ba725 "github.com/microsoftgraph/msgraph-sdk-go/scopedrolememberships/item"
    ie655917eacb97375016094a32b9a7b7696961d6c916c67c848b8c96ffa7e456e "github.com/microsoftgraph/msgraph-sdk-go/places/item"
    ie8e8e503fdb5d4623a3d8c75460ed99df6e7de79a047d61b15655f23eb0794f9 "github.com/microsoftgraph/msgraph-sdk-go/groupsettingtemplates/item"
    iebc0e64fadb20869bf2f248e5faa74af9d045c37a2822fb75e314761ad44656d "github.com/microsoftgraph/msgraph-sdk-go/oauth2permissiongrants/item"
    iec09f6187bc7a92a35b70b7fc70909dda436df66ea66bc31a862c86f0819cc15 "github.com/microsoftgraph/msgraph-sdk-go/directoryobjects/item"
    ieffc66507ab5f28c86663605f795e5c0be2a4353a94b34e8db2ddb67b0d285cf "github.com/microsoftgraph/msgraph-sdk-go/authenticationmethodconfigurations/item"
    if1f8863d252ff8f609844d73316e2ccaa8caf94c5c2e03b8a0aa91bcdc37a4bc "github.com/microsoftgraph/msgraph-sdk-go/groupsettings/item"
    if405c95e51d6685837bc60276ac44a0be46f00a5930cc59ce198c3a5119099a0 "github.com/microsoftgraph/msgraph-sdk-go/subscriptions/item"
    if697b88a41912c7fd65ee1e2a7160ff70a9f5fd5c48778b62d0d0ef4bc46fdb9 "github.com/microsoftgraph/msgraph-sdk-go/groups/item"
)

// GraphBaseServiceClient the main entry point of the SDK, exposes the configuration and the fluent API.
type GraphBaseServiceClient struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// Admin provides operations to manage the admin singleton.
func (m *GraphBaseServiceClient) Admin()(*i7c9d1b36ac198368c1d8bed014b43e2a518b170ee45bf02c8bbe64544a50539a.AdminRequestBuilder) {
    return i7c9d1b36ac198368c1d8bed014b43e2a518b170ee45bf02c8bbe64544a50539a.NewAdminRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AgreementAcceptances provides operations to manage the collection of agreementAcceptance entities.
func (m *GraphBaseServiceClient) AgreementAcceptances()(*i3e9b5129e2bb8b32b0374f7afe2536be6674d73df6c41d7c529f5a5432c4e0aa.AgreementAcceptancesRequestBuilder) {
    return i3e9b5129e2bb8b32b0374f7afe2536be6674d73df6c41d7c529f5a5432c4e0aa.NewAgreementAcceptancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AgreementAcceptancesById provides operations to manage the collection of agreementAcceptance entities.
func (m *GraphBaseServiceClient) AgreementAcceptancesById(id string)(*i16322b5a209cc0bd8b212f5c6ff6e00e3c536005637b4e1beb5a7076653b687a.AgreementAcceptanceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["agreementAcceptance%2Did"] = id
    }
    return i16322b5a209cc0bd8b212f5c6ff6e00e3c536005637b4e1beb5a7076653b687a.NewAgreementAcceptanceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Agreements provides operations to manage the collection of agreement entities.
func (m *GraphBaseServiceClient) Agreements()(*ieaa2790c8b7fa361674e69e4a385e279c8c641adf79d86e5b0ca566591a507e8.AgreementsRequestBuilder) {
    return ieaa2790c8b7fa361674e69e4a385e279c8c641adf79d86e5b0ca566591a507e8.NewAgreementsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AgreementsById provides operations to manage the collection of agreement entities.
func (m *GraphBaseServiceClient) AgreementsById(id string)(*i0832d6fdb5fdd25a7a1f3a0a99c90d6ef5c44fdc1f1e1f72f7d777d32cf5ea19.AgreementItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["agreement%2Did"] = id
    }
    return i0832d6fdb5fdd25a7a1f3a0a99c90d6ef5c44fdc1f1e1f72f7d777d32cf5ea19.NewAgreementItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AppCatalogs provides operations to manage the appCatalogs singleton.
func (m *GraphBaseServiceClient) AppCatalogs()(*i7d140130aac6882792a019b5ebe51fe8d69dfd63ec213c2e3cd98282ce2d0428.AppCatalogsRequestBuilder) {
    return i7d140130aac6882792a019b5ebe51fe8d69dfd63ec213c2e3cd98282ce2d0428.NewAppCatalogsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Applications provides operations to manage the collection of application entities.
func (m *GraphBaseServiceClient) Applications()(*i93194122344a685a2f9264205dc6d89a5ba39afdcea57fd0ade8f54b6f137c02.ApplicationsRequestBuilder) {
    return i93194122344a685a2f9264205dc6d89a5ba39afdcea57fd0ade8f54b6f137c02.NewApplicationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ApplicationsById provides operations to manage the collection of application entities.
func (m *GraphBaseServiceClient) ApplicationsById(id string)(*i1cb1089a14118ff1e1e1bd2d5f465b91b163bfc288e9bb57dc21502014e6b0c1.ApplicationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["application%2Did"] = id
    }
    return i1cb1089a14118ff1e1e1bd2d5f465b91b163bfc288e9bb57dc21502014e6b0c1.NewApplicationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ApplicationTemplates provides operations to manage the collection of applicationTemplate entities.
func (m *GraphBaseServiceClient) ApplicationTemplates()(*i5310ba7d4cfddbf5de4c1be94a30f9ca8c747c30a87e76587ce88d1cbfff01b4.ApplicationTemplatesRequestBuilder) {
    return i5310ba7d4cfddbf5de4c1be94a30f9ca8c747c30a87e76587ce88d1cbfff01b4.NewApplicationTemplatesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ApplicationTemplatesById provides operations to manage the collection of applicationTemplate entities.
func (m *GraphBaseServiceClient) ApplicationTemplatesById(id string)(*i2cc346a33133d41934cfa6e862c4eb0d4cf1dc36485198c479852b282338a897.ApplicationTemplateItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["applicationTemplate%2Did"] = id
    }
    return i2cc346a33133d41934cfa6e862c4eb0d4cf1dc36485198c479852b282338a897.NewApplicationTemplateItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AuditLogs provides operations to manage the auditLogRoot singleton.
func (m *GraphBaseServiceClient) AuditLogs()(*i20b08d3949f1191430a14a315e0758a1f131dc59bbdc93e654f1dd447a6af14c.AuditLogsRequestBuilder) {
    return i20b08d3949f1191430a14a315e0758a1f131dc59bbdc93e654f1dd447a6af14c.NewAuditLogsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AuthenticationMethodConfigurations provides operations to manage the collection of authenticationMethodConfiguration entities.
func (m *GraphBaseServiceClient) AuthenticationMethodConfigurations()(*i1dc06c4b7f499cb445a6c55e466abd6d7466bb35a2683c675909db23c57898e7.AuthenticationMethodConfigurationsRequestBuilder) {
    return i1dc06c4b7f499cb445a6c55e466abd6d7466bb35a2683c675909db23c57898e7.NewAuthenticationMethodConfigurationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AuthenticationMethodConfigurationsById provides operations to manage the collection of authenticationMethodConfiguration entities.
func (m *GraphBaseServiceClient) AuthenticationMethodConfigurationsById(id string)(*ieffc66507ab5f28c86663605f795e5c0be2a4353a94b34e8db2ddb67b0d285cf.AuthenticationMethodConfigurationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["authenticationMethodConfiguration%2Did"] = id
    }
    return ieffc66507ab5f28c86663605f795e5c0be2a4353a94b34e8db2ddb67b0d285cf.NewAuthenticationMethodConfigurationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AuthenticationMethodsPolicy provides operations to manage the authenticationMethodsPolicy singleton.
func (m *GraphBaseServiceClient) AuthenticationMethodsPolicy()(*id81f15a01b3ceaefa8b1b55f4ee944912f2179aafc4d873f0a2eaf0853eeccd0.AuthenticationMethodsPolicyRequestBuilder) {
    return id81f15a01b3ceaefa8b1b55f4ee944912f2179aafc4d873f0a2eaf0853eeccd0.NewAuthenticationMethodsPolicyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Branding provides operations to manage the organizationalBranding singleton.
func (m *GraphBaseServiceClient) Branding()(*if398f5c2f1cb53106e045240edd469d82f1854899fd95cfdf8f559b19375750c.BrandingRequestBuilder) {
    return if398f5c2f1cb53106e045240edd469d82f1854899fd95cfdf8f559b19375750c.NewBrandingRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CertificateBasedAuthConfiguration provides operations to manage the collection of certificateBasedAuthConfiguration entities.
func (m *GraphBaseServiceClient) CertificateBasedAuthConfiguration()(*i535d6c02ba98f73ff3a8c1c12a035ba5de51606f93aa2c0babdfed56fe505550.CertificateBasedAuthConfigurationRequestBuilder) {
    return i535d6c02ba98f73ff3a8c1c12a035ba5de51606f93aa2c0babdfed56fe505550.NewCertificateBasedAuthConfigurationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CertificateBasedAuthConfigurationById provides operations to manage the collection of certificateBasedAuthConfiguration entities.
func (m *GraphBaseServiceClient) CertificateBasedAuthConfigurationById(id string)(*i861ab10f93223993f014e54921d0feda5bc5dc9a8996dbbbb728e672d3e8162e.CertificateBasedAuthConfigurationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["certificateBasedAuthConfiguration%2Did"] = id
    }
    return i861ab10f93223993f014e54921d0feda5bc5dc9a8996dbbbb728e672d3e8162e.NewCertificateBasedAuthConfigurationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Chats provides operations to manage the collection of chat entities.
func (m *GraphBaseServiceClient) Chats()(*ibaef614e7692eebc6aaa8080b8ac29169fdf539f24925bc1de4465a3fcdac177.ChatsRequestBuilder) {
    return ibaef614e7692eebc6aaa8080b8ac29169fdf539f24925bc1de4465a3fcdac177.NewChatsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChatsById provides operations to manage the collection of chat entities.
func (m *GraphBaseServiceClient) ChatsById(id string)(*i4c038322c9990a1737934ea712b7d72201832cc50fbe4b6e382968270c3611c3.ChatItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["chat%2Did"] = id
    }
    return i4c038322c9990a1737934ea712b7d72201832cc50fbe4b6e382968270c3611c3.NewChatItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Communications provides operations to manage the cloudCommunications singleton.
func (m *GraphBaseServiceClient) Communications()(*i51b9802eedc1a25686534d117657be902df58c07e90ac6ea84501100998084d9.CommunicationsRequestBuilder) {
    return i51b9802eedc1a25686534d117657be902df58c07e90ac6ea84501100998084d9.NewCommunicationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Compliance provides operations to manage the compliance singleton.
func (m *GraphBaseServiceClient) Compliance()(*ia4b736f581ceef30e9ef8cebd9a6c2b932628e087982ff3dd2c9a0f1a920a918.ComplianceRequestBuilder) {
    return ia4b736f581ceef30e9ef8cebd9a6c2b932628e087982ff3dd2c9a0f1a920a918.NewComplianceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Connections provides operations to manage the collection of externalConnection entities.
func (m *GraphBaseServiceClient) Connections()(*icabdee72951e77325f237b36d388a199c87e65f67652b6bb85723aba847d7e83.ConnectionsRequestBuilder) {
    return icabdee72951e77325f237b36d388a199c87e65f67652b6bb85723aba847d7e83.NewConnectionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ConnectionsById provides operations to manage the collection of externalConnection entities.
func (m *GraphBaseServiceClient) ConnectionsById(id string)(*i728b903f252feeed28263ff4e0da95a1d0f7c831116db07abb00de41db959892.ExternalConnectionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["externalConnection%2Did"] = id
    }
    return i728b903f252feeed28263ff4e0da95a1d0f7c831116db07abb00de41db959892.NewExternalConnectionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// NewGraphBaseServiceClient instantiates a new GraphBaseServiceClient and sets the default values.
func NewGraphBaseServiceClient(requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*GraphBaseServiceClient) {
    m := &GraphBaseServiceClient{
    }
    m.pathParameters = make(map[string]string);
    m.urlTemplate = "{+baseurl}";
    m.requestAdapter = requestAdapter;
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RegisterDefaultSerializer(func() i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriterFactory { return i25911dc319edd61cbac496af7eab5ef20b6069a42515e22ec6a9bc97bf598488.NewJsonSerializationWriterFactory() })
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RegisterDefaultSerializer(func() i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriterFactory { return i7294a22093d408fdca300f11b81a887d89c47b764af06c8b803e2323973fdb83.NewTextSerializationWriterFactory() })
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RegisterDefaultDeserializer(func() i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNodeFactory { return i25911dc319edd61cbac496af7eab5ef20b6069a42515e22ec6a9bc97bf598488.NewJsonParseNodeFactory() })
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RegisterDefaultDeserializer(func() i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNodeFactory { return i7294a22093d408fdca300f11b81a887d89c47b764af06c8b803e2323973fdb83.NewTextParseNodeFactory() })
    if m.requestAdapter.GetBaseUrl() == "" {
        m.requestAdapter.SetBaseUrl("https://graph.microsoft.com/v1.0")
    }
    return m
}
// Contacts provides operations to manage the collection of orgContact entities.
func (m *GraphBaseServiceClient) Contacts()(*if51cca2652371587dbc02e65260e291435a6a8f7f2ffb419f26c3b9d2a033f57.ContactsRequestBuilder) {
    return if51cca2652371587dbc02e65260e291435a6a8f7f2ffb419f26c3b9d2a033f57.NewContactsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ContactsById provides operations to manage the collection of orgContact entities.
func (m *GraphBaseServiceClient) ContactsById(id string)(*i3ea73b2c3d03b959a8d3906c2b1b073951feb694dd6984ab1eea4e6b8c1d0858.OrgContactItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["orgContact%2Did"] = id
    }
    return i3ea73b2c3d03b959a8d3906c2b1b073951feb694dd6984ab1eea4e6b8c1d0858.NewOrgContactItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Contracts provides operations to manage the collection of contract entities.
func (m *GraphBaseServiceClient) Contracts()(*ie3631868038c44f490dbc03525ac7249d0523c29cc45cbb25b2aebcf470d6c0c.ContractsRequestBuilder) {
    return ie3631868038c44f490dbc03525ac7249d0523c29cc45cbb25b2aebcf470d6c0c.NewContractsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ContractsById provides operations to manage the collection of contract entities.
func (m *GraphBaseServiceClient) ContractsById(id string)(*i86bb3d05e1a6bbdb496bd3c65829f1a6eb272be42e9ac6060f873dfbd921e4ea.ContractItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["contract%2Did"] = id
    }
    return i86bb3d05e1a6bbdb496bd3c65829f1a6eb272be42e9ac6060f873dfbd921e4ea.NewContractItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DataPolicyOperations provides operations to manage the collection of dataPolicyOperation entities.
func (m *GraphBaseServiceClient) DataPolicyOperations()(*ib33fc5e9889e020c0c572578957f59819123a589c61fd7f3eb37eb7958b525ee.DataPolicyOperationsRequestBuilder) {
    return ib33fc5e9889e020c0c572578957f59819123a589c61fd7f3eb37eb7958b525ee.NewDataPolicyOperationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DataPolicyOperationsById provides operations to manage the collection of dataPolicyOperation entities.
func (m *GraphBaseServiceClient) DataPolicyOperationsById(id string)(*i206b88e4abcec8b25d993b477b59473f3e9420358bcd878d07c5ed2b531ccf4c.DataPolicyOperationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["dataPolicyOperation%2Did"] = id
    }
    return i206b88e4abcec8b25d993b477b59473f3e9420358bcd878d07c5ed2b531ccf4c.NewDataPolicyOperationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DeviceAppManagement provides operations to manage the deviceAppManagement singleton.
func (m *GraphBaseServiceClient) DeviceAppManagement()(*i638650494f9db477daff56d31ff923f5c100f72df0257ed7fa5c222cb1a77a94.DeviceAppManagementRequestBuilder) {
    return i638650494f9db477daff56d31ff923f5c100f72df0257ed7fa5c222cb1a77a94.NewDeviceAppManagementRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceManagement provides operations to manage the deviceManagement singleton.
func (m *GraphBaseServiceClient) DeviceManagement()(*i738daeb889f22c1e163aee5a37a094b55b1d815dc76d4802d64e4e1b2e44206c.DeviceManagementRequestBuilder) {
    return i738daeb889f22c1e163aee5a37a094b55b1d815dc76d4802d64e4e1b2e44206c.NewDeviceManagementRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Devices provides operations to manage the collection of device entities.
func (m *GraphBaseServiceClient) Devices()(*i4c91eeb51f03f9d59a342065f7c6ee027ad1fe84ada6b1946b8162c5ae146cfb.DevicesRequestBuilder) {
    return i4c91eeb51f03f9d59a342065f7c6ee027ad1fe84ada6b1946b8162c5ae146cfb.NewDevicesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DevicesById provides operations to manage the collection of device entities.
func (m *GraphBaseServiceClient) DevicesById(id string)(*ib6d66da0f7d4860b7205f5fdb1200fc9000adb4fbc853a2f05f70c644580220f.DeviceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["device%2Did"] = id
    }
    return ib6d66da0f7d4860b7205f5fdb1200fc9000adb4fbc853a2f05f70c644580220f.NewDeviceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Directory provides operations to manage the directory singleton.
func (m *GraphBaseServiceClient) Directory()(*ib14d748b564c787931c10f1c7ba6856eeddea29a5b9e5c5c27eb1224ff65e5c4.DirectoryRequestBuilder) {
    return ib14d748b564c787931c10f1c7ba6856eeddea29a5b9e5c5c27eb1224ff65e5c4.NewDirectoryRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DirectoryObjects provides operations to manage the collection of directoryObject entities.
func (m *GraphBaseServiceClient) DirectoryObjects()(*i185698f71f6301975f0627ee999e6e91920d8fa9c00bdef3487b9f349e2df04e.DirectoryObjectsRequestBuilder) {
    return i185698f71f6301975f0627ee999e6e91920d8fa9c00bdef3487b9f349e2df04e.NewDirectoryObjectsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DirectoryObjectsById provides operations to manage the collection of directoryObject entities.
func (m *GraphBaseServiceClient) DirectoryObjectsById(id string)(*iec09f6187bc7a92a35b70b7fc70909dda436df66ea66bc31a862c86f0819cc15.DirectoryObjectItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryObject%2Did"] = id
    }
    return iec09f6187bc7a92a35b70b7fc70909dda436df66ea66bc31a862c86f0819cc15.NewDirectoryObjectItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DirectoryRoles provides operations to manage the collection of directoryRole entities.
func (m *GraphBaseServiceClient) DirectoryRoles()(*id4615a956cb1e7edabf8f5a4bc131d1ceca9a13d0f79ae0e122997452a9a0a4e.DirectoryRolesRequestBuilder) {
    return id4615a956cb1e7edabf8f5a4bc131d1ceca9a13d0f79ae0e122997452a9a0a4e.NewDirectoryRolesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DirectoryRolesById provides operations to manage the collection of directoryRole entities.
func (m *GraphBaseServiceClient) DirectoryRolesById(id string)(*i960f074bae2d1f849aec23c162b7be41055a1485f8efd075e3c89717cc4ac8f5.DirectoryRoleItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryRole%2Did"] = id
    }
    return i960f074bae2d1f849aec23c162b7be41055a1485f8efd075e3c89717cc4ac8f5.NewDirectoryRoleItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DirectoryRoleTemplates provides operations to manage the collection of directoryRoleTemplate entities.
func (m *GraphBaseServiceClient) DirectoryRoleTemplates()(*i2a252d42835bdab6d88bf938595da6cf029001f9ca970d6f599cecf0ca27f8e5.DirectoryRoleTemplatesRequestBuilder) {
    return i2a252d42835bdab6d88bf938595da6cf029001f9ca970d6f599cecf0ca27f8e5.NewDirectoryRoleTemplatesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DirectoryRoleTemplatesById provides operations to manage the collection of directoryRoleTemplate entities.
func (m *GraphBaseServiceClient) DirectoryRoleTemplatesById(id string)(*icf62d3bb4e29c8437041430705851ef556cb3cf91d77df26e8eaf92a32e9814e.DirectoryRoleTemplateItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryRoleTemplate%2Did"] = id
    }
    return icf62d3bb4e29c8437041430705851ef556cb3cf91d77df26e8eaf92a32e9814e.NewDirectoryRoleTemplateItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DomainDnsRecords provides operations to manage the collection of domainDnsRecord entities.
func (m *GraphBaseServiceClient) DomainDnsRecords()(*iaca6694a878291d0e4021155b406c19d3080cdfc382b456e43c71264d4d9e519.DomainDnsRecordsRequestBuilder) {
    return iaca6694a878291d0e4021155b406c19d3080cdfc382b456e43c71264d4d9e519.NewDomainDnsRecordsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DomainDnsRecordsById provides operations to manage the collection of domainDnsRecord entities.
func (m *GraphBaseServiceClient) DomainDnsRecordsById(id string)(*i55644ead1eb4b291341a5935abcd8425f7456cccabc4594ef4aee967d51fc863.DomainDnsRecordItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["domainDnsRecord%2Did"] = id
    }
    return i55644ead1eb4b291341a5935abcd8425f7456cccabc4594ef4aee967d51fc863.NewDomainDnsRecordItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Domains provides operations to manage the collection of domain entities.
func (m *GraphBaseServiceClient) Domains()(*i957076b10ba162b23efec7b94dd26b84c6475d285449c1cbc9c5b85910d36a12.DomainsRequestBuilder) {
    return i957076b10ba162b23efec7b94dd26b84c6475d285449c1cbc9c5b85910d36a12.NewDomainsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DomainsById provides operations to manage the collection of domain entities.
func (m *GraphBaseServiceClient) DomainsById(id string)(*i164d67321f322030ad927126612d90d5880d461e3357bd32611832679c00aff2.DomainItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["domain%2Did"] = id
    }
    return i164d67321f322030ad927126612d90d5880d461e3357bd32611832679c00aff2.NewDomainItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Drive provides operations to manage the drive singleton.
func (m *GraphBaseServiceClient) Drive()(*i926bd489c52af20f44aacc8a450bb0a062290f1d1e44c2fe78d6cc1595c12524.DriveRequestBuilder) {
    return i926bd489c52af20f44aacc8a450bb0a062290f1d1e44c2fe78d6cc1595c12524.NewDriveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Drives provides operations to manage the collection of drive entities.
func (m *GraphBaseServiceClient) Drives()(*iefc72d8a17962d4db125c50866617eaa15d662c6e3fb13735d477380dcc0dbe3.DrivesRequestBuilder) {
    return iefc72d8a17962d4db125c50866617eaa15d662c6e3fb13735d477380dcc0dbe3.NewDrivesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DrivesById provides operations to manage the collection of drive entities.
func (m *GraphBaseServiceClient) DrivesById(id string)(*ic73c2557206a32cd6d6e58acca928163e176fe1cf7382bdae4f55d28ff56e345.DriveItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["drive%2Did"] = id
    }
    return ic73c2557206a32cd6d6e58acca928163e176fe1cf7382bdae4f55d28ff56e345.NewDriveItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Education provides operations to manage the educationRoot singleton.
func (m *GraphBaseServiceClient) Education()(*i1be0f1b1da466bc62355d411ef490acbd8dc0ec5ca4d3448c7eb73e5caffafc3.EducationRequestBuilder) {
    return i1be0f1b1da466bc62355d411ef490acbd8dc0ec5ca4d3448c7eb73e5caffafc3.NewEducationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// External provides operations to manage the external singleton.
func (m *GraphBaseServiceClient) External()(*ib3217193884e00033cb8182cac52178dfa3b20ce9c4eb48e37a6217882d956ae.ExternalRequestBuilder) {
    return ib3217193884e00033cb8182cac52178dfa3b20ce9c4eb48e37a6217882d956ae.NewExternalRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GroupLifecyclePolicies provides operations to manage the collection of groupLifecyclePolicy entities.
func (m *GraphBaseServiceClient) GroupLifecyclePolicies()(*i1d6652ecc686b20c37a9a3448b26db8187e284e1a4017cab8876b02b97557436.GroupLifecyclePoliciesRequestBuilder) {
    return i1d6652ecc686b20c37a9a3448b26db8187e284e1a4017cab8876b02b97557436.NewGroupLifecyclePoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GroupLifecyclePoliciesById provides operations to manage the collection of groupLifecyclePolicy entities.
func (m *GraphBaseServiceClient) GroupLifecyclePoliciesById(id string)(*icd43f65a6cebd6c2685c244d7f46f49a951d8647717cded612ba79705e6aa7c7.GroupLifecyclePolicyItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["groupLifecyclePolicy%2Did"] = id
    }
    return icd43f65a6cebd6c2685c244d7f46f49a951d8647717cded612ba79705e6aa7c7.NewGroupLifecyclePolicyItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Groups provides operations to manage the collection of group entities.
func (m *GraphBaseServiceClient) Groups()(*ia6e876e3ed2d92c29c13dbc8c37513bc38d0d5f05ab9321e43a25ff336912a2d.GroupsRequestBuilder) {
    return ia6e876e3ed2d92c29c13dbc8c37513bc38d0d5f05ab9321e43a25ff336912a2d.NewGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GroupsById provides operations to manage the collection of group entities.
func (m *GraphBaseServiceClient) GroupsById(id string)(*if697b88a41912c7fd65ee1e2a7160ff70a9f5fd5c48778b62d0d0ef4bc46fdb9.GroupItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["group%2Did"] = id
    }
    return if697b88a41912c7fd65ee1e2a7160ff70a9f5fd5c48778b62d0d0ef4bc46fdb9.NewGroupItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// GroupSettings provides operations to manage the collection of groupSetting entities.
func (m *GraphBaseServiceClient) GroupSettings()(*i4794c103c0d044c27a3ca3af0a0e498e93a9863420c1a4e7a29ef37590053c7b.GroupSettingsRequestBuilder) {
    return i4794c103c0d044c27a3ca3af0a0e498e93a9863420c1a4e7a29ef37590053c7b.NewGroupSettingsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GroupSettingsById provides operations to manage the collection of groupSetting entities.
func (m *GraphBaseServiceClient) GroupSettingsById(id string)(*if1f8863d252ff8f609844d73316e2ccaa8caf94c5c2e03b8a0aa91bcdc37a4bc.GroupSettingItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["groupSetting%2Did"] = id
    }
    return if1f8863d252ff8f609844d73316e2ccaa8caf94c5c2e03b8a0aa91bcdc37a4bc.NewGroupSettingItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// GroupSettingTemplates provides operations to manage the collection of groupSettingTemplate entities.
func (m *GraphBaseServiceClient) GroupSettingTemplates()(*id2ac823944414906187dbe4e6ca3b5e46886b9db738d2c1c27de6df8b1bebd61.GroupSettingTemplatesRequestBuilder) {
    return id2ac823944414906187dbe4e6ca3b5e46886b9db738d2c1c27de6df8b1bebd61.NewGroupSettingTemplatesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GroupSettingTemplatesById provides operations to manage the collection of groupSettingTemplate entities.
func (m *GraphBaseServiceClient) GroupSettingTemplatesById(id string)(*ie8e8e503fdb5d4623a3d8c75460ed99df6e7de79a047d61b15655f23eb0794f9.GroupSettingTemplateItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["groupSettingTemplate%2Did"] = id
    }
    return ie8e8e503fdb5d4623a3d8c75460ed99df6e7de79a047d61b15655f23eb0794f9.NewGroupSettingTemplateItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Identity provides operations to manage the identityContainer singleton.
func (m *GraphBaseServiceClient) Identity()(*i79ca23a9ac0659e1330dd29e049fe157787d5af6695ead2ff8263396db68d027.IdentityRequestBuilder) {
    return i79ca23a9ac0659e1330dd29e049fe157787d5af6695ead2ff8263396db68d027.NewIdentityRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// IdentityGovernance provides operations to manage the identityGovernance singleton.
func (m *GraphBaseServiceClient) IdentityGovernance()(*i32d45c1243c349600fbe53b2f9641bb59857a3326037587cbe4e347b46ad207e.IdentityGovernanceRequestBuilder) {
    return i32d45c1243c349600fbe53b2f9641bb59857a3326037587cbe4e347b46ad207e.NewIdentityGovernanceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// IdentityProtection provides operations to manage the identityProtectionRoot singleton.
func (m *GraphBaseServiceClient) IdentityProtection()(*i663c30678b300c2c4b619c4964b4326e471e4da61a44d7c39f752349da7a468e.IdentityProtectionRequestBuilder) {
    return i663c30678b300c2c4b619c4964b4326e471e4da61a44d7c39f752349da7a468e.NewIdentityProtectionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// IdentityProviders provides operations to manage the collection of identityProvider entities.
func (m *GraphBaseServiceClient) IdentityProviders()(*i62c2771f3f3a1e5e085aedcde54473e9f043cc57b9ce4dd88980a77aca7a5a10.IdentityProvidersRequestBuilder) {
    return i62c2771f3f3a1e5e085aedcde54473e9f043cc57b9ce4dd88980a77aca7a5a10.NewIdentityProvidersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// IdentityProvidersById provides operations to manage the collection of identityProvider entities.
func (m *GraphBaseServiceClient) IdentityProvidersById(id string)(*i292f565d97fd74d4ba58ae7a10fd84e8982b27c0e3ba8747d8fd7f8b61482271.IdentityProviderItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["identityProvider%2Did"] = id
    }
    return i292f565d97fd74d4ba58ae7a10fd84e8982b27c0e3ba8747d8fd7f8b61482271.NewIdentityProviderItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// InformationProtection provides operations to manage the informationProtection singleton.
func (m *GraphBaseServiceClient) InformationProtection()(*ib68fa8e66bda853b3a33c491e8a66ca665897dab129192b2c97289266c4a1415.InformationProtectionRequestBuilder) {
    return ib68fa8e66bda853b3a33c491e8a66ca665897dab129192b2c97289266c4a1415.NewInformationProtectionRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Invitations provides operations to manage the collection of invitation entities.
func (m *GraphBaseServiceClient) Invitations()(*ic5e701d75e87f15ce153687b00984a314f7eeea8cfdc77cd9ad648e5ccbc7fbd.InvitationsRequestBuilder) {
    return ic5e701d75e87f15ce153687b00984a314f7eeea8cfdc77cd9ad648e5ccbc7fbd.NewInvitationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InvitationsById provides operations to manage the collection of invitation entities.
func (m *GraphBaseServiceClient) InvitationsById(id string)(*i62e63d0f24e0caea5c3b2202a0ee0923bdf496f3e118faa4ee49895e02eecfde.InvitationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["invitation%2Did"] = id
    }
    return i62e63d0f24e0caea5c3b2202a0ee0923bdf496f3e118faa4ee49895e02eecfde.NewInvitationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Localizations provides operations to manage the collection of organizationalBrandingLocalization entities.
func (m *GraphBaseServiceClient) Localizations()(*i61686672307beee899fe5a14188df42982da47730f55a14800b102cd10ab2d72.LocalizationsRequestBuilder) {
    return i61686672307beee899fe5a14188df42982da47730f55a14800b102cd10ab2d72.NewLocalizationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// LocalizationsById provides operations to manage the collection of organizationalBrandingLocalization entities.
func (m *GraphBaseServiceClient) LocalizationsById(id string)(*i24a463345a902b042d0fa0b40e03cab9b230ab6328a113241a23ab1b81c0bcd1.OrganizationalBrandingLocalizationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["organizationalBrandingLocalization%2Did"] = id
    }
    return i24a463345a902b042d0fa0b40e03cab9b230ab6328a113241a23ab1b81c0bcd1.NewOrganizationalBrandingLocalizationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Me provides operations to manage the user singleton.
func (m *GraphBaseServiceClient) Me()(*i71117da372286e863c042a526ec1361696ab14b838a5b77db5bc54386d436543.MeRequestBuilder) {
    return i71117da372286e863c042a526ec1361696ab14b838a5b77db5bc54386d436543.NewMeRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Oauth2PermissionGrants provides operations to manage the collection of oAuth2PermissionGrant entities.
func (m *GraphBaseServiceClient) Oauth2PermissionGrants()(*i86cada4d4a5f2f8a9d1e7a85eacd70a661ea7b20d2737008c0719e95b5be3e16.Oauth2PermissionGrantsRequestBuilder) {
    return i86cada4d4a5f2f8a9d1e7a85eacd70a661ea7b20d2737008c0719e95b5be3e16.NewOauth2PermissionGrantsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Oauth2PermissionGrantsById provides operations to manage the collection of oAuth2PermissionGrant entities.
func (m *GraphBaseServiceClient) Oauth2PermissionGrantsById(id string)(*iebc0e64fadb20869bf2f248e5faa74af9d045c37a2822fb75e314761ad44656d.OAuth2PermissionGrantItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["oAuth2PermissionGrant%2Did"] = id
    }
    return iebc0e64fadb20869bf2f248e5faa74af9d045c37a2822fb75e314761ad44656d.NewOAuth2PermissionGrantItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Organization provides operations to manage the collection of organization entities.
func (m *GraphBaseServiceClient) Organization()(*ic949a0bb5066d68760e8502a7f9db83f571d9e01e38fad4aadf7268188e52df0.OrganizationRequestBuilder) {
    return ic949a0bb5066d68760e8502a7f9db83f571d9e01e38fad4aadf7268188e52df0.NewOrganizationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OrganizationById provides operations to manage the collection of organization entities.
func (m *GraphBaseServiceClient) OrganizationById(id string)(*ia5422e2deae41871358311d10a7b0d4a60e914828d2fe80bf0a1bd96c1ff2a2f.OrganizationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["organization%2Did"] = id
    }
    return ia5422e2deae41871358311d10a7b0d4a60e914828d2fe80bf0a1bd96c1ff2a2f.NewOrganizationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// PermissionGrants provides operations to manage the collection of resourceSpecificPermissionGrant entities.
func (m *GraphBaseServiceClient) PermissionGrants()(*i1a1369b1521a8ac4885166fd68eae4247248a891006fea464d2eea2a271b2cdb.PermissionGrantsRequestBuilder) {
    return i1a1369b1521a8ac4885166fd68eae4247248a891006fea464d2eea2a271b2cdb.NewPermissionGrantsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PermissionGrantsById provides operations to manage the collection of resourceSpecificPermissionGrant entities.
func (m *GraphBaseServiceClient) PermissionGrantsById(id string)(*i23bab38fb8688d4bab0b6ffc533eb085d40e58af49a27ab228a8d1ad3e0ab203.ResourceSpecificPermissionGrantItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["resourceSpecificPermissionGrant%2Did"] = id
    }
    return i23bab38fb8688d4bab0b6ffc533eb085d40e58af49a27ab228a8d1ad3e0ab203.NewResourceSpecificPermissionGrantItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Places provides operations to manage the collection of place entities.
func (m *GraphBaseServiceClient) Places()(*ia3e0f7c2d21d5c73ecb8a7552177d0fe444ae0522290dd1c4b5559e449b118af.PlacesRequestBuilder) {
    return ia3e0f7c2d21d5c73ecb8a7552177d0fe444ae0522290dd1c4b5559e449b118af.NewPlacesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PlacesById provides operations to manage the collection of place entities.
func (m *GraphBaseServiceClient) PlacesById(id string)(*ie655917eacb97375016094a32b9a7b7696961d6c916c67c848b8c96ffa7e456e.PlaceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["place%2Did"] = id
    }
    return ie655917eacb97375016094a32b9a7b7696961d6c916c67c848b8c96ffa7e456e.NewPlaceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Planner provides operations to manage the planner singleton.
func (m *GraphBaseServiceClient) Planner()(*ie66b913c1bc1c536bc8db5d185910e9318f621374e016f95e36e9d59b7127f63.PlannerRequestBuilder) {
    return ie66b913c1bc1c536bc8db5d185910e9318f621374e016f95e36e9d59b7127f63.NewPlannerRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Policies provides operations to manage the policyRoot singleton.
func (m *GraphBaseServiceClient) Policies()(*i4ac7f0a844871066493521918f268cafe2a25c71c28a98221ea3f22d5153090f.PoliciesRequestBuilder) {
    return i4ac7f0a844871066493521918f268cafe2a25c71c28a98221ea3f22d5153090f.NewPoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Print provides operations to manage the print singleton.
func (m *GraphBaseServiceClient) Print()(*i9429d7aae2f5c1dabbecc9411e8ad2b733d29338bc0c0436eeccc94605c461b7.PrintRequestBuilder) {
    return i9429d7aae2f5c1dabbecc9411e8ad2b733d29338bc0c0436eeccc94605c461b7.NewPrintRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Privacy provides operations to manage the privacy singleton.
func (m *GraphBaseServiceClient) Privacy()(*i58857a108d6e260e56ef0dd7e783668388f113eb436006780703ac59f0abb3b1.PrivacyRequestBuilder) {
    return i58857a108d6e260e56ef0dd7e783668388f113eb436006780703ac59f0abb3b1.NewPrivacyRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Reports provides operations to manage the reportRoot singleton.
func (m *GraphBaseServiceClient) Reports()(*i6bf2d83eea06710580ad0d54b886ac4e14cbab0d1d84937f340f02b99f8f5738.ReportsRequestBuilder) {
    return i6bf2d83eea06710580ad0d54b886ac4e14cbab0d1d84937f340f02b99f8f5738.NewReportsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RoleManagement provides operations to manage the roleManagement singleton.
func (m *GraphBaseServiceClient) RoleManagement()(*if5555fa41b6637688bcf8c25c62a041258f4dc6eacb38ad42d91c66f222ee182.RoleManagementRequestBuilder) {
    return if5555fa41b6637688bcf8c25c62a041258f4dc6eacb38ad42d91c66f222ee182.NewRoleManagementRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SchemaExtensions provides operations to manage the collection of schemaExtension entities.
func (m *GraphBaseServiceClient) SchemaExtensions()(*i4a624e38d68c2a9fc4db1ea915bcaffde116f967f58ec2c99e2ea8bbff3690e1.SchemaExtensionsRequestBuilder) {
    return i4a624e38d68c2a9fc4db1ea915bcaffde116f967f58ec2c99e2ea8bbff3690e1.NewSchemaExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SchemaExtensionsById provides operations to manage the collection of schemaExtension entities.
func (m *GraphBaseServiceClient) SchemaExtensionsById(id string)(*i31ea098bea3bcc5d0f40d4471abad6488910c0e4872783682f25fb2bdbcffdb6.SchemaExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["schemaExtension%2Did"] = id
    }
    return i31ea098bea3bcc5d0f40d4471abad6488910c0e4872783682f25fb2bdbcffdb6.NewSchemaExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ScopedRoleMemberships provides operations to manage the collection of scopedRoleMembership entities.
func (m *GraphBaseServiceClient) ScopedRoleMemberships()(*ice10f31b9db59ba91184d2b882172edb754f885050cf0830aa2b7c8ff880556b.ScopedRoleMembershipsRequestBuilder) {
    return ice10f31b9db59ba91184d2b882172edb754f885050cf0830aa2b7c8ff880556b.NewScopedRoleMembershipsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ScopedRoleMembershipsById provides operations to manage the collection of scopedRoleMembership entities.
func (m *GraphBaseServiceClient) ScopedRoleMembershipsById(id string)(*id5e9a05bae8f5cd30c57fd87690f009f004424eafeb45208f44e7ed46f8ba725.ScopedRoleMembershipItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["scopedRoleMembership%2Did"] = id
    }
    return id5e9a05bae8f5cd30c57fd87690f009f004424eafeb45208f44e7ed46f8ba725.NewScopedRoleMembershipItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Search provides operations to manage the searchEntity singleton.
func (m *GraphBaseServiceClient) Search()(*i286f3babd79fe9ec3b0f52b6ed5910842c0adaeff02be1843d0e01c56d9ba6d9.SearchRequestBuilder) {
    return i286f3babd79fe9ec3b0f52b6ed5910842c0adaeff02be1843d0e01c56d9ba6d9.NewSearchRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Security provides operations to manage the security singleton.
func (m *GraphBaseServiceClient) Security()(*i35d7bbcc8f7e8b8e9525ea0ee5b3c51c3a1a58f9ed512b727d181bfcd08eb032.SecurityRequestBuilder) {
    return i35d7bbcc8f7e8b8e9525ea0ee5b3c51c3a1a58f9ed512b727d181bfcd08eb032.NewSecurityRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipals provides operations to manage the collection of servicePrincipal entities.
func (m *GraphBaseServiceClient) ServicePrincipals()(*i009f47bbce65ccdb7303730eed71e6bab3ae2f8e4e918bc9e94341d28624af97.ServicePrincipalsRequestBuilder) {
    return i009f47bbce65ccdb7303730eed71e6bab3ae2f8e4e918bc9e94341d28624af97.NewServicePrincipalsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ServicePrincipalsById provides operations to manage the collection of servicePrincipal entities.
func (m *GraphBaseServiceClient) ServicePrincipalsById(id string)(*i120b7d5508b5c9e9e49c562cc2c54282d0cac65c8ec72e8928f45cc697956915.ServicePrincipalItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["servicePrincipal%2Did"] = id
    }
    return i120b7d5508b5c9e9e49c562cc2c54282d0cac65c8ec72e8928f45cc697956915.NewServicePrincipalItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Shares provides operations to manage the collection of sharedDriveItem entities.
func (m *GraphBaseServiceClient) Shares()(*i07d47a144340607d6d6dbd93575e531530e4f1cc6091c947ea0766f7951ffd34.SharesRequestBuilder) {
    return i07d47a144340607d6d6dbd93575e531530e4f1cc6091c947ea0766f7951ffd34.NewSharesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SharesById provides operations to manage the collection of sharedDriveItem entities.
func (m *GraphBaseServiceClient) SharesById(id string)(*i1a7ecef3f68d1e375afc82029319b84e73d0840cf0fa231262dfbe877c5b0f30.SharedDriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["sharedDriveItem%2Did"] = id
    }
    return i1a7ecef3f68d1e375afc82029319b84e73d0840cf0fa231262dfbe877c5b0f30.NewSharedDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Sites provides operations to manage the collection of site entities.
func (m *GraphBaseServiceClient) Sites()(*if5372351befdb652f617b1ee71fbf092fa8dd2a161ba9c021bc265628b6ea82b.SitesRequestBuilder) {
    return if5372351befdb652f617b1ee71fbf092fa8dd2a161ba9c021bc265628b6ea82b.NewSitesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SitesById provides operations to manage the collection of site entities.
func (m *GraphBaseServiceClient) SitesById(id string)(*i7bd49b9f046ac5c4108447b2999dd223be9f057c9910bdc197bcba51be5ac641.SiteItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["site%2Did"] = id
    }
    return i7bd49b9f046ac5c4108447b2999dd223be9f057c9910bdc197bcba51be5ac641.NewSiteItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Solutions provides operations to manage the solutionsRoot singleton.
func (m *GraphBaseServiceClient) Solutions()(*i0906e75d8a44bf92212e084e1d2f62d03887dcec6a5c8535e92ccc04c1e5fdec.SolutionsRequestBuilder) {
    return i0906e75d8a44bf92212e084e1d2f62d03887dcec6a5c8535e92ccc04c1e5fdec.NewSolutionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubscribedSkus provides operations to manage the collection of subscribedSku entities.
func (m *GraphBaseServiceClient) SubscribedSkus()(*id007bc768abbff1131aab64890cdcd0411159a946e9df27140c5f7cf8f249647.SubscribedSkusRequestBuilder) {
    return id007bc768abbff1131aab64890cdcd0411159a946e9df27140c5f7cf8f249647.NewSubscribedSkusRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubscribedSkusById provides operations to manage the collection of subscribedSku entities.
func (m *GraphBaseServiceClient) SubscribedSkusById(id string)(*i838f22f4ef3405018163f34ed56f20e99e7bb1a5bd0bfd7009a72d054a09a36c.SubscribedSkuItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["subscribedSku%2Did"] = id
    }
    return i838f22f4ef3405018163f34ed56f20e99e7bb1a5bd0bfd7009a72d054a09a36c.NewSubscribedSkuItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Subscriptions provides operations to manage the collection of subscription entities.
func (m *GraphBaseServiceClient) Subscriptions()(*idb8230b65f4a369c23b4d9b41ebe568c657c92f8f77fe36d16d64528b3a317a3.SubscriptionsRequestBuilder) {
    return idb8230b65f4a369c23b4d9b41ebe568c657c92f8f77fe36d16d64528b3a317a3.NewSubscriptionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// SubscriptionsById provides operations to manage the collection of subscription entities.
func (m *GraphBaseServiceClient) SubscriptionsById(id string)(*if405c95e51d6685837bc60276ac44a0be46f00a5930cc59ce198c3a5119099a0.SubscriptionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["subscription%2Did"] = id
    }
    return if405c95e51d6685837bc60276ac44a0be46f00a5930cc59ce198c3a5119099a0.NewSubscriptionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Teams provides operations to manage the collection of team entities.
func (m *GraphBaseServiceClient) Teams()(*if39bc788926a05e976b265ecfc616408ca12af399df9ce3a2bb348fe89708057.TeamsRequestBuilder) {
    return if39bc788926a05e976b265ecfc616408ca12af399df9ce3a2bb348fe89708057.NewTeamsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TeamsById provides operations to manage the collection of team entities.
func (m *GraphBaseServiceClient) TeamsById(id string)(*i4b66adb64e72d30468c3d30826f1df6386f1703ce53cf4a46d2a6cbbfd88b747.TeamItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["team%2Did"] = id
    }
    return i4b66adb64e72d30468c3d30826f1df6386f1703ce53cf4a46d2a6cbbfd88b747.NewTeamItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TeamsTemplates provides operations to manage the collection of teamsTemplate entities.
func (m *GraphBaseServiceClient) TeamsTemplates()(*i1b75be7b5675627960b4672ab148be21ff379d5cbc0e62f6bc5b97d54464f8b5.TeamsTemplatesRequestBuilder) {
    return i1b75be7b5675627960b4672ab148be21ff379d5cbc0e62f6bc5b97d54464f8b5.NewTeamsTemplatesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TeamsTemplatesById provides operations to manage the collection of teamsTemplate entities.
func (m *GraphBaseServiceClient) TeamsTemplatesById(id string)(*i2a87b2a2ef0c6367c3abc2e909d9176a42e3124d1896753ad0a25d9a6a881c32.TeamsTemplateItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["teamsTemplate%2Did"] = id
    }
    return i2a87b2a2ef0c6367c3abc2e909d9176a42e3124d1896753ad0a25d9a6a881c32.NewTeamsTemplateItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Teamwork provides operations to manage the teamwork singleton.
func (m *GraphBaseServiceClient) Teamwork()(*ie05ac24b652f7d895cca374316c093c4ca40dd2df0f1518c465233d6432b1ef9.TeamworkRequestBuilder) {
    return ie05ac24b652f7d895cca374316c093c4ca40dd2df0f1518c465233d6432b1ef9.NewTeamworkRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Users provides operations to manage the collection of user entities.
func (m *GraphBaseServiceClient) Users()(*if6ffd1464db2d9c22e351b03e4c00ebd24a5353cd70ffb7f56cfad1c3ceec329.UsersRequestBuilder) {
    return if6ffd1464db2d9c22e351b03e4c00ebd24a5353cd70ffb7f56cfad1c3ceec329.NewUsersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// UsersById provides operations to manage the collection of user entities.
func (m *GraphBaseServiceClient) UsersById(id string)(*i009581390843c78f63b06f9dcefeeb5ef2a124a2ac1dcbad3adbe4d0d5e650af.UserItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["user%2Did"] = id
    }
    return i009581390843c78f63b06f9dcefeeb5ef2a124a2ac1dcbad3adbe4d0d5e650af.NewUserItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Workbooks provides operations to manage the collection of driveItem entities.
func (m *GraphBaseServiceClient) Workbooks()(*i088cccaaaff488138f258ec600ca804c5dc9548772ebe52bf3cb7a3eaf4b9fdf.WorkbooksRequestBuilder) {
    return i088cccaaaff488138f258ec600ca804c5dc9548772ebe52bf3cb7a3eaf4b9fdf.NewWorkbooksRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// WorkbooksById provides operations to manage the collection of driveItem entities.
func (m *GraphBaseServiceClient) WorkbooksById(id string)(*i31a01e7736f53f5c2e0acd0c456ed6dcc6065ca402aaf31b6a1526bd3bf57394.DriveItemItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["driveItem%2Did"] = id
    }
    return i31a01e7736f53f5c2e0acd0c456ed6dcc6065ca402aaf31b6a1526bd3bf57394.NewDriveItemItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
