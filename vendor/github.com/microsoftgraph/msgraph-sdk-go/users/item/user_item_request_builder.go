package item

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    i005f0821c1e57467ed9e31e0f65ebf8795241942e8afc0c6c9cc866dc3b70de2 "github.com/microsoftgraph/msgraph-sdk-go/users/item/owneddevices"
    i065d0a279baf0cb315c26855c2ba7a5c0701939d8cae543a4fb9b7bb2c7cbda5 "github.com/microsoftgraph/msgraph-sdk-go/users/item/licensedetails"
    i0a25da984a2612b28dedf3c1e18748fd27dd92433244095b38a29e6ec842f917 "github.com/microsoftgraph/msgraph-sdk-go/users/item/checkmemberobjects"
    i0a4c7a84985f2276fcfbcf69e2fe37b80ef87c2622cc8763371022616974902d "github.com/microsoftgraph/msgraph-sdk-go/users/item/presence"
    i0a6d7999baeed47d1aa02960cd9f3f120528ecae4240bf52300542a8dc9f978d "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages"
    i1090e84ad9b76c4c279467ebbc90ff155a540b45b098e680706d2aae300639eb "github.com/microsoftgraph/msgraph-sdk-go/users/item/revokesigninsessions"
    i11db6285564609d0b36b9032e0d8f6d257828b887e558375da4354121f40a3d0 "github.com/microsoftgraph/msgraph-sdk-go/users/item/transitivememberof"
    i152c9659f52d4d76fff6d586e3de05874eee179419588197c21a3adc1e3a6c7b "github.com/microsoftgraph/msgraph-sdk-go/users/item/ownedobjects"
    i16f13fb20dc21e9cbd93cbca432ac5cb871c47a97e5bebce67f5eab6fe8276e9 "github.com/microsoftgraph/msgraph-sdk-go/users/item/sendmail"
    i225d659ae554743be8acfa963ab6384d1f6e82911866fa2e3c01918b6a0637b1 "github.com/microsoftgraph/msgraph-sdk-go/users/item/settings"
    i23845c92400f42d2b469653fba5fe1555d60b50e4596b30d2d2b4c127b0ceecf "github.com/microsoftgraph/msgraph-sdk-go/users/item/getmanagedapppolicies"
    i2b4d48c4393811d79f4d4cc1f4addb3b49f58a53edcf7d1d7140e29a22ab7504 "github.com/microsoftgraph/msgraph-sdk-go/users/item/people"
    i2c0f098b49b5f2edba81f4095f0c699516263b7cd6ff316d31f469f15b8366b0 "github.com/microsoftgraph/msgraph-sdk-go/users/item/onlinemeetings"
    i2ec0496e560d7cecb96d82bcf32b4c61894a21cfd165c378253939195c546d27 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars"
    i3a4cc4e91af7383718bf555760e1dba5c1fc1de3509837d9074d3a3f0a8a0fc5 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drive"
    i3fd2d5130e27e5ed481f7a0565ee17271d52e0ea15f0dec8ae436a3780d1d8d2 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manager"
    i418d77dc601ca4dc2c568aef8702ce31b38463f4e1701c6c788adcbeb32ca3e8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/outlook"
    i4a474d43148dd96e52e36b3ed7c52a7d75f46d95adf4fa0ce96ace2509d33744 "github.com/microsoftgraph/msgraph-sdk-go/users/item/planner"
    i4af8697d288ced3f62e5f7467084c929a355a7f4c2140a3f51072b00f73663c1 "github.com/microsoftgraph/msgraph-sdk-go/users/item/wipemanagedappregistrationsbydevicetag"
    i4cb75aef1c5547ab5604eb9b1721354472df7129995c61c066ac036f69c410c1 "github.com/microsoftgraph/msgraph-sdk-go/users/item/createdobjects"
    i54c28e8ec50e0dd4eeee8073ed79cecee11ba12e3813ec711d8127af3d4b994c "github.com/microsoftgraph/msgraph-sdk-go/users/item/getmembergroups"
    i61d84fd6521641b94f790414a5b9abccf870ea66012c270c06f3ea7df770aad8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/todo"
    i631d94a7300d4142a2d7a4c4d0d424e98d3fb058f106e6567a8ac6476f5956b6 "github.com/microsoftgraph/msgraph-sdk-go/users/item/registereddevices"
    i64f47462143d9c5d6613e409adca3c0fd24631de22e340968821ffdff467e657 "github.com/microsoftgraph/msgraph-sdk-go/users/item/inferenceclassification"
    i6d278e5cb359938f2362c090cc322d0c80867015903a39849ea478a6cda51689 "github.com/microsoftgraph/msgraph-sdk-go/users/item/removealldevicesfrommanagement"
    i6d5ec618068b4f9940966e3235a6b56547b161354e35247fc9843cd2fc68c280 "github.com/microsoftgraph/msgraph-sdk-go/users/item/getmailtips"
    i6fea53c850cb6db38ea2669df26d5c9986df82967fc9dcbd1415ef4c37ba3909 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives"
    i7434742f328fb34181df9069343c7bf11d068d61f4481ab2eaea71f3bc398a8c "github.com/microsoftgraph/msgraph-sdk-go/users/item/activities"
    i79cfaa5f7496a3faeffd9cee89b05478df04194a9746fc39a36637ad15ede76b "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendar"
    i7f31fc68d418ebfd772842294b63ed5b4dafc3cab499c7bf25b7634a901656bb "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams"
    i8acf48036fa30b3733f0a3efffd923a945db096ce1c005ede1236ec3db8f202b "github.com/microsoftgraph/msgraph-sdk-go/users/item/approleassignments"
    i8b015d6ac3b9f998d8da5aa503d326b952acdc90c4c90a03d48fdec51c4588f1 "github.com/microsoftgraph/msgraph-sdk-go/users/item/devicemanagementtroubleshootingevents"
    i8deae6de37c6b836f34cd5c6b8fab91a73043c26ee346f7410a8dcf88ea64c99 "github.com/microsoftgraph/msgraph-sdk-go/users/item/insights"
    i8fdd3e26eeb6da18c0b9d487ac43983355a059891790fba9c732ef59d5c4f21b "github.com/microsoftgraph/msgraph-sdk-go/users/item/findmeetingtimes"
    i9164851f0a06aba10b88da1bbecb98e3600356d368c9db499b8931c4f96f1219 "github.com/microsoftgraph/msgraph-sdk-go/users/item/contacts"
    i91af54e73e80a81d2f98c6f72fbb6c81539e67ba0480227ae8e2e5f2b8f467fa "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices"
    i91bcb6ebe27791a826a51ecff16cc35355876378f254587360afeb0054d6f984 "github.com/microsoftgraph/msgraph-sdk-go/users/item/checkmembergroups"
    i927d56fd90d9c38a0e88ef5c682feb77cc18445308d44276e0f1e095247b3634 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups"
    i94ff5ee566023b18b8b2ac8c1da92081f3bf4520205017599c00aed9faff56de "github.com/microsoftgraph/msgraph-sdk-go/users/item/teamwork"
    i9f87171771d54744395d1043beb42c14430d7cb4f1557826ebef7c5b2584f7cd "github.com/microsoftgraph/msgraph-sdk-go/users/item/restore"
    ia6a0f4f70d60f824df18f07f8918c6e9ac043f588293118e5d91eb921d5dffbf "github.com/microsoftgraph/msgraph-sdk-go/users/item/memberof"
    ia99f649e03b3020696f5c2c11f42a9ad7f094cb80add889863fc7947d3bf9ec7 "github.com/microsoftgraph/msgraph-sdk-go/users/item/scopedrolememberof"
    ib8e8ec0d09197f1fd57167b1b67da9bcc229fa24977921ce6406f24aa89a6884 "github.com/microsoftgraph/msgraph-sdk-go/users/item/onenote"
    ib946e3869f342130dddab8c76c89f264a76e899f07657349241151b4e41b4286 "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders"
    ic31c6bdf1febcc1619240151a993f96448ebdf3c6f91fe0243fbe55f91e7f620 "github.com/microsoftgraph/msgraph-sdk-go/users/item/contactfolders"
    icf27f39169518adb5f7aa45892636b30915e5166590a8d55b561e502c71dc5fe "github.com/microsoftgraph/msgraph-sdk-go/users/item/changepassword"
    icf56cc7fcccc9803d5655a975f93bee1c76d14d27e2e180c40ef7bb1c3c18a68 "github.com/microsoftgraph/msgraph-sdk-go/users/item/photo"
    id7592a93402c0d67508da968e58d07109c9c2561ae838fbeb1cccb27ef558ecb "github.com/microsoftgraph/msgraph-sdk-go/users/item/directreports"
    id9510e26530c8ba3b39834c060eb72230844f30169f4e2b2e6b683eb86727fa0 "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats"
    idc22dd282127655d723afbe2b6bcc58326fe6d9a8b1bcdee2586b77850b235b0 "github.com/microsoftgraph/msgraph-sdk-go/users/item/agreementacceptances"
    iddc156b6ba287ee33d5592797e8da64087fd068cf179bd892944253be67ac174 "github.com/microsoftgraph/msgraph-sdk-go/users/item/reminderviewwithstartdatetimewithenddatetime"
    ide0a33943859ca09b2034a1804f60d32a108dc09bf7012d68e5d3c63ac617ddb "github.com/microsoftgraph/msgraph-sdk-go/users/item/getmemberobjects"
    ide2d4144e5fd4f180ce2f7d5791881ad155e598ebbe5bda78dbfc5f9cd8aa8d0 "github.com/microsoftgraph/msgraph-sdk-go/users/item/authentication"
    idea2a9c81deaaa3e33a2bcef2386387935820f4f571d71af1fc2daa24d8fa032 "github.com/microsoftgraph/msgraph-sdk-go/users/item/photos"
    idfea0d624d269e60a96eeda8949da3675c3ceb484b3fb16f118382ba0a7d3583 "github.com/microsoftgraph/msgraph-sdk-go/users/item/assignlicense"
    iea339692a4ece4755a03d26004b71e7bbfe3ad28cf8036bdc4b972512eb00359 "github.com/microsoftgraph/msgraph-sdk-go/users/item/oauth2permissiongrants"
    ieb873ea17846be13ef65aaf00090f933a16727baa2e8063b39c2d86ae07e1ad0 "github.com/microsoftgraph/msgraph-sdk-go/users/item/getmanagedappdiagnosticstatuses"
    iedc2060488fc2f6e7d120554e9e853117ec2aab876ef464c9ed3da4696541a32 "github.com/microsoftgraph/msgraph-sdk-go/users/item/exportpersonaldata"
    if0c7df43a05145fae76c1ab18483108fb4043f7b6a0b01f47772fd66022550ee "github.com/microsoftgraph/msgraph-sdk-go/users/item/extensions"
    if12f3532feee7012e3a4dc1e3a2444ca8fab8f29e8128baafe0ef00f6ed6015c "github.com/microsoftgraph/msgraph-sdk-go/users/item/managedappregistrations"
    if19172f362f2e1d4456e041244db6eae76a05acf4f16807f05721a5ef0a0f4be "github.com/microsoftgraph/msgraph-sdk-go/users/item/followedsites"
    ifd940c65ea13f1270f41f753ee5e7761f074915bd6a6c64b2adb77f70aea8c32 "github.com/microsoftgraph/msgraph-sdk-go/users/item/reprocesslicenseassignment"
    ifd9b701a6f88aa459780777b5ac521fe3494f01b9276c6884b86827296582e54 "github.com/microsoftgraph/msgraph-sdk-go/users/item/events"
    ifdf31129bec8038d82607a6d973830e0b090629373d5812152e5046514a4bfb0 "github.com/microsoftgraph/msgraph-sdk-go/users/item/translateexchangeids"
    ifec6c3ff45ea43d5b46ee636228680b7e4b3bfb9eba0437938b28dad5cadcd06 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview"
    i03f3617a0bebd4a3033a4ba021c004c146b37855bfeed1510a905b18d838b51e "github.com/microsoftgraph/msgraph-sdk-go/users/item/licensedetails/item"
    i077531fcb7f8b5dcbcd7beae97fc12c083a5c8a8e949fc043865f6a338711296 "github.com/microsoftgraph/msgraph-sdk-go/users/item/directreports/item"
    i086f912938edbc37b845789cda2e3e84c9a478d06edbee6d535f8e143d4c70ba "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendars/item"
    i0b4dff6b82c5c98b58b9ff11aead66a7e49d0eac6b18b1490fd01bf9597cf36c "github.com/microsoftgraph/msgraph-sdk-go/users/item/joinedteams/item"
    i0be6421da86ae0eb5550fa2bf3e3bb1ee5c6cb82a7b6539fce14162126c5d994 "github.com/microsoftgraph/msgraph-sdk-go/users/item/devicemanagementtroubleshootingevents/item"
    i0c6b3d37dc64615f253a6fd9a58f59f536aa93a57f99bb5f9d4210ae83aa5112 "github.com/microsoftgraph/msgraph-sdk-go/users/item/manageddevices/item"
    i129160b87cb78d35acdd1a84c7328c6872f51bf17cb02a211ba957a5431db64b "github.com/microsoftgraph/msgraph-sdk-go/users/item/extensions/item"
    i2083055461886ca3b56f8adc26419c83b6c45f08c3a0eb4fcdfa94b7e99c9bef "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendargroups/item"
    i3468dc539adf033fc51d31c6b79eb77caea5661d9461c49781e068c23230b9fc "github.com/microsoftgraph/msgraph-sdk-go/users/item/registereddevices/item"
    i3ecd250b638acba95f41b3db07eaa507842a51f8e0cd3640cb9b75f0052ceff6 "github.com/microsoftgraph/msgraph-sdk-go/users/item/onlinemeetings/item"
    i49037195fa7f3333c97e06f8131490f42ed26328e1c7d9d27cba659e29c03a56 "github.com/microsoftgraph/msgraph-sdk-go/users/item/oauth2permissiongrants/item"
    i4d1b56e56093574975d1ab3f4a52af1d6db962510708352cb01245aa891a6ef5 "github.com/microsoftgraph/msgraph-sdk-go/users/item/chats/item"
    i4fd4546e4cf3f7b84a0adef775c8a8fc52aee1e036e23fc206091f5b7b54638e "github.com/microsoftgraph/msgraph-sdk-go/users/item/agreementacceptances/item"
    i5835fd9ad2a434b5df38fb1f3e5d0ec77c2f4ba2d0a321c93550f076a0281720 "github.com/microsoftgraph/msgraph-sdk-go/users/item/managedappregistrations/item"
    i5903dc47170759db544ed9d0a0c895f70016bf360d13d02d94dd87cf4f13cd5d "github.com/microsoftgraph/msgraph-sdk-go/users/item/memberof/item"
    i5cf53dd103ab02b1568c2f84b946a72977e41e0de97287af6722ccafaa880d49 "github.com/microsoftgraph/msgraph-sdk-go/users/item/calendarview/item"
    i5f2b6c8998555350079828979727c205406cf7155e907008bec42440c3e9724a "github.com/microsoftgraph/msgraph-sdk-go/users/item/photos/item"
    i755269849b1d98c9e29fa525a149a4aede0223fcbc1aefdb163bc3d1fc4aa32e "github.com/microsoftgraph/msgraph-sdk-go/users/item/owneddevices/item"
    i764f7513465c0251a3785c28eac14dab3642e586e0fad67528f13704ed1cbb64 "github.com/microsoftgraph/msgraph-sdk-go/users/item/drives/item"
    i8465a2d5f26666d08a0baebc755ef48f1f0dc18551f929c58c09d0110510a317 "github.com/microsoftgraph/msgraph-sdk-go/users/item/people/item"
    i88e2c30b78a734ede88b94c583c45b04ce07232920ad9d7b2366b4278bfae360 "github.com/microsoftgraph/msgraph-sdk-go/users/item/contacts/item"
    i983bdc3dce2d947aec82e1e9dd5edd887e8d861d589287d050f96af67abb1e1f "github.com/microsoftgraph/msgraph-sdk-go/users/item/scopedrolememberof/item"
    ibbdc262c9fd17ec4015663fbf467ebc865f1a6a45daf1ffbaa73c799220a610b "github.com/microsoftgraph/msgraph-sdk-go/users/item/contactfolders/item"
    id1f5db1c2220a25bede6871151a00a067652bdb75e2db4d199ca6ffcdd70b342 "github.com/microsoftgraph/msgraph-sdk-go/users/item/mailfolders/item"
    id307eb3180e14c2c4cec59b4b85564ae282ebaa15d7f7735716952f69e678279 "github.com/microsoftgraph/msgraph-sdk-go/users/item/approleassignments/item"
    id3793eb5f17f42ec0817eacc0b999985de40ba38c196a38c7bd6b0c736c5345f "github.com/microsoftgraph/msgraph-sdk-go/users/item/createdobjects/item"
    id761928b5250e3d39114f0e3b14e8ea044e71e53b81f38c27c8e5836522b3e1b "github.com/microsoftgraph/msgraph-sdk-go/users/item/followedsites/item"
    ie9198f71b9fb192ecc84221b88715d18fc05b44b87d4b863d0e98c3c6b88c246 "github.com/microsoftgraph/msgraph-sdk-go/users/item/transitivememberof/item"
    if1c7e6e4f2864e681dc58961189ef8bfe8cee67a94c24f3b9f32c1ffafdd2aa2 "github.com/microsoftgraph/msgraph-sdk-go/users/item/ownedobjects/item"
    if2fde6f9b7cf9a66adc51b4c84ec599a0497b6792379909cd577d454bd9d0414 "github.com/microsoftgraph/msgraph-sdk-go/users/item/activities/item"
    if33563b345c020375260161400a6b809f00973f17d5ab870445f3c67ca7323b8 "github.com/microsoftgraph/msgraph-sdk-go/users/item/messages/item"
    ifdc39e295bcc2e19fecc18f2ae737ddb8fe9018d8753e81ef0c2de8bb37638da "github.com/microsoftgraph/msgraph-sdk-go/users/item/events/item"
)

// UserItemRequestBuilder provides operations to manage the collection of user entities.
type UserItemRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// UserItemRequestBuilderDeleteRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type UserItemRequestBuilderDeleteRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// UserItemRequestBuilderGetQueryParameters retrieve the properties and relationships of user object.
type UserItemRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// UserItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type UserItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *UserItemRequestBuilderGetQueryParameters
}
// UserItemRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type UserItemRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Activities provides operations to manage the activities property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Activities()(*i7434742f328fb34181df9069343c7bf11d068d61f4481ab2eaea71f3bc398a8c.ActivitiesRequestBuilder) {
    return i7434742f328fb34181df9069343c7bf11d068d61f4481ab2eaea71f3bc398a8c.NewActivitiesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ActivitiesById provides operations to manage the activities property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) ActivitiesById(id string)(*if2fde6f9b7cf9a66adc51b4c84ec599a0497b6792379909cd577d454bd9d0414.UserActivityItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["userActivity%2Did"] = id
    }
    return if2fde6f9b7cf9a66adc51b4c84ec599a0497b6792379909cd577d454bd9d0414.NewUserActivityItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AgreementAcceptances provides operations to manage the agreementAcceptances property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) AgreementAcceptances()(*idc22dd282127655d723afbe2b6bcc58326fe6d9a8b1bcdee2586b77850b235b0.AgreementAcceptancesRequestBuilder) {
    return idc22dd282127655d723afbe2b6bcc58326fe6d9a8b1bcdee2586b77850b235b0.NewAgreementAcceptancesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AgreementAcceptancesById provides operations to manage the agreementAcceptances property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) AgreementAcceptancesById(id string)(*i4fd4546e4cf3f7b84a0adef775c8a8fc52aee1e036e23fc206091f5b7b54638e.AgreementAcceptanceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["agreementAcceptance%2Did"] = id
    }
    return i4fd4546e4cf3f7b84a0adef775c8a8fc52aee1e036e23fc206091f5b7b54638e.NewAgreementAcceptanceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AppRoleAssignments provides operations to manage the appRoleAssignments property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) AppRoleAssignments()(*i8acf48036fa30b3733f0a3efffd923a945db096ce1c005ede1236ec3db8f202b.AppRoleAssignmentsRequestBuilder) {
    return i8acf48036fa30b3733f0a3efffd923a945db096ce1c005ede1236ec3db8f202b.NewAppRoleAssignmentsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// AppRoleAssignmentsById provides operations to manage the appRoleAssignments property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) AppRoleAssignmentsById(id string)(*id307eb3180e14c2c4cec59b4b85564ae282ebaa15d7f7735716952f69e678279.AppRoleAssignmentItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["appRoleAssignment%2Did"] = id
    }
    return id307eb3180e14c2c4cec59b4b85564ae282ebaa15d7f7735716952f69e678279.NewAppRoleAssignmentItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// AssignLicense provides operations to call the assignLicense method.
func (m *UserItemRequestBuilder) AssignLicense()(*idfea0d624d269e60a96eeda8949da3675c3ceb484b3fb16f118382ba0a7d3583.AssignLicenseRequestBuilder) {
    return idfea0d624d269e60a96eeda8949da3675c3ceb484b3fb16f118382ba0a7d3583.NewAssignLicenseRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Authentication provides operations to manage the authentication property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Authentication()(*ide2d4144e5fd4f180ce2f7d5791881ad155e598ebbe5bda78dbfc5f9cd8aa8d0.AuthenticationRequestBuilder) {
    return ide2d4144e5fd4f180ce2f7d5791881ad155e598ebbe5bda78dbfc5f9cd8aa8d0.NewAuthenticationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Calendar provides operations to manage the calendar property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Calendar()(*i79cfaa5f7496a3faeffd9cee89b05478df04194a9746fc39a36637ad15ede76b.CalendarRequestBuilder) {
    return i79cfaa5f7496a3faeffd9cee89b05478df04194a9746fc39a36637ad15ede76b.NewCalendarRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CalendarGroups provides operations to manage the calendarGroups property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) CalendarGroups()(*i927d56fd90d9c38a0e88ef5c682feb77cc18445308d44276e0f1e095247b3634.CalendarGroupsRequestBuilder) {
    return i927d56fd90d9c38a0e88ef5c682feb77cc18445308d44276e0f1e095247b3634.NewCalendarGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CalendarGroupsById provides operations to manage the calendarGroups property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) CalendarGroupsById(id string)(*i2083055461886ca3b56f8adc26419c83b6c45f08c3a0eb4fcdfa94b7e99c9bef.CalendarGroupItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["calendarGroup%2Did"] = id
    }
    return i2083055461886ca3b56f8adc26419c83b6c45f08c3a0eb4fcdfa94b7e99c9bef.NewCalendarGroupItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Calendars provides operations to manage the calendars property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Calendars()(*i2ec0496e560d7cecb96d82bcf32b4c61894a21cfd165c378253939195c546d27.CalendarsRequestBuilder) {
    return i2ec0496e560d7cecb96d82bcf32b4c61894a21cfd165c378253939195c546d27.NewCalendarsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CalendarsById provides operations to manage the calendars property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) CalendarsById(id string)(*i086f912938edbc37b845789cda2e3e84c9a478d06edbee6d535f8e143d4c70ba.CalendarItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["calendar%2Did"] = id
    }
    return i086f912938edbc37b845789cda2e3e84c9a478d06edbee6d535f8e143d4c70ba.NewCalendarItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// CalendarView provides operations to manage the calendarView property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) CalendarView()(*ifec6c3ff45ea43d5b46ee636228680b7e4b3bfb9eba0437938b28dad5cadcd06.CalendarViewRequestBuilder) {
    return ifec6c3ff45ea43d5b46ee636228680b7e4b3bfb9eba0437938b28dad5cadcd06.NewCalendarViewRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CalendarViewById provides operations to manage the calendarView property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) CalendarViewById(id string)(*i5cf53dd103ab02b1568c2f84b946a72977e41e0de97287af6722ccafaa880d49.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did"] = id
    }
    return i5cf53dd103ab02b1568c2f84b946a72977e41e0de97287af6722ccafaa880d49.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ChangePassword provides operations to call the changePassword method.
func (m *UserItemRequestBuilder) ChangePassword()(*icf27f39169518adb5f7aa45892636b30915e5166590a8d55b561e502c71dc5fe.ChangePasswordRequestBuilder) {
    return icf27f39169518adb5f7aa45892636b30915e5166590a8d55b561e502c71dc5fe.NewChangePasswordRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Chats provides operations to manage the chats property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Chats()(*id9510e26530c8ba3b39834c060eb72230844f30169f4e2b2e6b683eb86727fa0.ChatsRequestBuilder) {
    return id9510e26530c8ba3b39834c060eb72230844f30169f4e2b2e6b683eb86727fa0.NewChatsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ChatsById provides operations to manage the chats property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) ChatsById(id string)(*i4d1b56e56093574975d1ab3f4a52af1d6db962510708352cb01245aa891a6ef5.ChatItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["chat%2Did"] = id
    }
    return i4d1b56e56093574975d1ab3f4a52af1d6db962510708352cb01245aa891a6ef5.NewChatItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// CheckMemberGroups provides operations to call the checkMemberGroups method.
func (m *UserItemRequestBuilder) CheckMemberGroups()(*i91bcb6ebe27791a826a51ecff16cc35355876378f254587360afeb0054d6f984.CheckMemberGroupsRequestBuilder) {
    return i91bcb6ebe27791a826a51ecff16cc35355876378f254587360afeb0054d6f984.NewCheckMemberGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CheckMemberObjects provides operations to call the checkMemberObjects method.
func (m *UserItemRequestBuilder) CheckMemberObjects()(*i0a25da984a2612b28dedf3c1e18748fd27dd92433244095b38a29e6ec842f917.CheckMemberObjectsRequestBuilder) {
    return i0a25da984a2612b28dedf3c1e18748fd27dd92433244095b38a29e6ec842f917.NewCheckMemberObjectsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// NewUserItemRequestBuilderInternal instantiates a new UserItemRequestBuilder and sets the default values.
func NewUserItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*UserItemRequestBuilder) {
    m := &UserItemRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/users/{user%2Did}{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewUserItemRequestBuilder instantiates a new UserItemRequestBuilder and sets the default values.
func NewUserItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*UserItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewUserItemRequestBuilderInternal(urlParams, requestAdapter)
}
// ContactFolders provides operations to manage the contactFolders property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) ContactFolders()(*ic31c6bdf1febcc1619240151a993f96448ebdf3c6f91fe0243fbe55f91e7f620.ContactFoldersRequestBuilder) {
    return ic31c6bdf1febcc1619240151a993f96448ebdf3c6f91fe0243fbe55f91e7f620.NewContactFoldersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ContactFoldersById provides operations to manage the contactFolders property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) ContactFoldersById(id string)(*ibbdc262c9fd17ec4015663fbf467ebc865f1a6a45daf1ffbaa73c799220a610b.ContactFolderItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["contactFolder%2Did"] = id
    }
    return ibbdc262c9fd17ec4015663fbf467ebc865f1a6a45daf1ffbaa73c799220a610b.NewContactFolderItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Contacts provides operations to manage the contacts property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Contacts()(*i9164851f0a06aba10b88da1bbecb98e3600356d368c9db499b8931c4f96f1219.ContactsRequestBuilder) {
    return i9164851f0a06aba10b88da1bbecb98e3600356d368c9db499b8931c4f96f1219.NewContactsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ContactsById provides operations to manage the contacts property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) ContactsById(id string)(*i88e2c30b78a734ede88b94c583c45b04ce07232920ad9d7b2366b4278bfae360.ContactItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["contact%2Did"] = id
    }
    return i88e2c30b78a734ede88b94c583c45b04ce07232920ad9d7b2366b4278bfae360.NewContactItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// CreateDeleteRequestInformation delete user.   When deleted, user resources are moved to a temporary container and can be restored within 30 days.  After that time, they are permanently deleted.  To learn more, see deletedItems.
func (m *UserItemRequestBuilder) CreateDeleteRequestInformation(ctx context.Context, requestConfiguration *UserItemRequestBuilderDeleteRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.DELETE
    if requestConfiguration != nil {
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// CreatedObjects provides operations to manage the createdObjects property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) CreatedObjects()(*i4cb75aef1c5547ab5604eb9b1721354472df7129995c61c066ac036f69c410c1.CreatedObjectsRequestBuilder) {
    return i4cb75aef1c5547ab5604eb9b1721354472df7129995c61c066ac036f69c410c1.NewCreatedObjectsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// CreatedObjectsById provides operations to manage the createdObjects property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) CreatedObjectsById(id string)(*id3793eb5f17f42ec0817eacc0b999985de40ba38c196a38c7bd6b0c736c5345f.DirectoryObjectItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryObject%2Did"] = id
    }
    return id3793eb5f17f42ec0817eacc0b999985de40ba38c196a38c7bd6b0c736c5345f.NewDirectoryObjectItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// CreateGetRequestInformation retrieve the properties and relationships of user object.
func (m *UserItemRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *UserItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET
    requestInfo.Headers["Accept"] = "application/json"
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// CreatePatchRequestInformation update the properties of a user object. Not all properties can be updated by Member or Guest users with their default permissions without Administrator roles. Compare member and guest default permissions to see properties they can manage.
func (m *UserItemRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable, requestConfiguration *UserItemRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.PATCH
    requestInfo.Headers["Accept"] = "application/json"
    requestInfo.SetContentFromParsable(ctx, m.requestAdapter, "application/json", body)
    if requestConfiguration != nil {
        requestInfo.AddRequestHeaders(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// Delete delete user.   When deleted, user resources are moved to a temporary container and can be restored within 30 days.  After that time, they are permanently deleted.  To learn more, see deletedItems.
func (m *UserItemRequestBuilder) Delete(ctx context.Context, requestConfiguration *UserItemRequestBuilderDeleteRequestConfiguration)(error) {
    requestInfo, err := m.CreateDeleteRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    err = m.requestAdapter.SendNoContentAsync(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// DeviceManagementTroubleshootingEvents provides operations to manage the deviceManagementTroubleshootingEvents property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) DeviceManagementTroubleshootingEvents()(*i8b015d6ac3b9f998d8da5aa503d326b952acdc90c4c90a03d48fdec51c4588f1.DeviceManagementTroubleshootingEventsRequestBuilder) {
    return i8b015d6ac3b9f998d8da5aa503d326b952acdc90c4c90a03d48fdec51c4588f1.NewDeviceManagementTroubleshootingEventsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceManagementTroubleshootingEventsById provides operations to manage the deviceManagementTroubleshootingEvents property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) DeviceManagementTroubleshootingEventsById(id string)(*i0be6421da86ae0eb5550fa2bf3e3bb1ee5c6cb82a7b6539fce14162126c5d994.DeviceManagementTroubleshootingEventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["deviceManagementTroubleshootingEvent%2Did"] = id
    }
    return i0be6421da86ae0eb5550fa2bf3e3bb1ee5c6cb82a7b6539fce14162126c5d994.NewDeviceManagementTroubleshootingEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DirectReports provides operations to manage the directReports property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) DirectReports()(*id7592a93402c0d67508da968e58d07109c9c2561ae838fbeb1cccb27ef558ecb.DirectReportsRequestBuilder) {
    return id7592a93402c0d67508da968e58d07109c9c2561ae838fbeb1cccb27ef558ecb.NewDirectReportsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DirectReportsById provides operations to manage the directReports property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) DirectReportsById(id string)(*i077531fcb7f8b5dcbcd7beae97fc12c083a5c8a8e949fc043865f6a338711296.DirectoryObjectItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryObject%2Did"] = id
    }
    return i077531fcb7f8b5dcbcd7beae97fc12c083a5c8a8e949fc043865f6a338711296.NewDirectoryObjectItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Drive provides operations to manage the drive property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Drive()(*i3a4cc4e91af7383718bf555760e1dba5c1fc1de3509837d9074d3a3f0a8a0fc5.DriveRequestBuilder) {
    return i3a4cc4e91af7383718bf555760e1dba5c1fc1de3509837d9074d3a3f0a8a0fc5.NewDriveRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Drives provides operations to manage the drives property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Drives()(*i6fea53c850cb6db38ea2669df26d5c9986df82967fc9dcbd1415ef4c37ba3909.DrivesRequestBuilder) {
    return i6fea53c850cb6db38ea2669df26d5c9986df82967fc9dcbd1415ef4c37ba3909.NewDrivesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DrivesById provides operations to manage the drives property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) DrivesById(id string)(*i764f7513465c0251a3785c28eac14dab3642e586e0fad67528f13704ed1cbb64.DriveItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["drive%2Did"] = id
    }
    return i764f7513465c0251a3785c28eac14dab3642e586e0fad67528f13704ed1cbb64.NewDriveItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Events provides operations to manage the events property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Events()(*ifd9b701a6f88aa459780777b5ac521fe3494f01b9276c6884b86827296582e54.EventsRequestBuilder) {
    return ifd9b701a6f88aa459780777b5ac521fe3494f01b9276c6884b86827296582e54.NewEventsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// EventsById provides operations to manage the events property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) EventsById(id string)(*ifdc39e295bcc2e19fecc18f2ae737ddb8fe9018d8753e81ef0c2de8bb37638da.EventItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["event%2Did"] = id
    }
    return ifdc39e295bcc2e19fecc18f2ae737ddb8fe9018d8753e81ef0c2de8bb37638da.NewEventItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ExportPersonalData provides operations to call the exportPersonalData method.
func (m *UserItemRequestBuilder) ExportPersonalData()(*iedc2060488fc2f6e7d120554e9e853117ec2aab876ef464c9ed3da4696541a32.ExportPersonalDataRequestBuilder) {
    return iedc2060488fc2f6e7d120554e9e853117ec2aab876ef464c9ed3da4696541a32.NewExportPersonalDataRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Extensions provides operations to manage the extensions property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Extensions()(*if0c7df43a05145fae76c1ab18483108fb4043f7b6a0b01f47772fd66022550ee.ExtensionsRequestBuilder) {
    return if0c7df43a05145fae76c1ab18483108fb4043f7b6a0b01f47772fd66022550ee.NewExtensionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ExtensionsById provides operations to manage the extensions property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) ExtensionsById(id string)(*i129160b87cb78d35acdd1a84c7328c6872f51bf17cb02a211ba957a5431db64b.ExtensionItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["extension%2Did"] = id
    }
    return i129160b87cb78d35acdd1a84c7328c6872f51bf17cb02a211ba957a5431db64b.NewExtensionItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// FindMeetingTimes provides operations to call the findMeetingTimes method.
func (m *UserItemRequestBuilder) FindMeetingTimes()(*i8fdd3e26eeb6da18c0b9d487ac43983355a059891790fba9c732ef59d5c4f21b.FindMeetingTimesRequestBuilder) {
    return i8fdd3e26eeb6da18c0b9d487ac43983355a059891790fba9c732ef59d5c4f21b.NewFindMeetingTimesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FollowedSites provides operations to manage the followedSites property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) FollowedSites()(*if19172f362f2e1d4456e041244db6eae76a05acf4f16807f05721a5ef0a0f4be.FollowedSitesRequestBuilder) {
    return if19172f362f2e1d4456e041244db6eae76a05acf4f16807f05721a5ef0a0f4be.NewFollowedSitesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// FollowedSitesById provides operations to manage the followedSites property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) FollowedSitesById(id string)(*id761928b5250e3d39114f0e3b14e8ea044e71e53b81f38c27c8e5836522b3e1b.SiteItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["site%2Did"] = id
    }
    return id761928b5250e3d39114f0e3b14e8ea044e71e53b81f38c27c8e5836522b3e1b.NewSiteItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Get retrieve the properties and relationships of user object.
func (m *UserItemRequestBuilder) Get(ctx context.Context, requestConfiguration *UserItemRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUserFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable), nil
}
// GetMailTips provides operations to call the getMailTips method.
func (m *UserItemRequestBuilder) GetMailTips()(*i6d5ec618068b4f9940966e3235a6b56547b161354e35247fc9843cd2fc68c280.GetMailTipsRequestBuilder) {
    return i6d5ec618068b4f9940966e3235a6b56547b161354e35247fc9843cd2fc68c280.NewGetMailTipsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetManagedAppDiagnosticStatuses provides operations to call the getManagedAppDiagnosticStatuses method.
func (m *UserItemRequestBuilder) GetManagedAppDiagnosticStatuses()(*ieb873ea17846be13ef65aaf00090f933a16727baa2e8063b39c2d86ae07e1ad0.GetManagedAppDiagnosticStatusesRequestBuilder) {
    return ieb873ea17846be13ef65aaf00090f933a16727baa2e8063b39c2d86ae07e1ad0.NewGetManagedAppDiagnosticStatusesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetManagedAppPolicies provides operations to call the getManagedAppPolicies method.
func (m *UserItemRequestBuilder) GetManagedAppPolicies()(*i23845c92400f42d2b469653fba5fe1555d60b50e4596b30d2d2b4c127b0ceecf.GetManagedAppPoliciesRequestBuilder) {
    return i23845c92400f42d2b469653fba5fe1555d60b50e4596b30d2d2b4c127b0ceecf.NewGetManagedAppPoliciesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetMemberGroups provides operations to call the getMemberGroups method.
func (m *UserItemRequestBuilder) GetMemberGroups()(*i54c28e8ec50e0dd4eeee8073ed79cecee11ba12e3813ec711d8127af3d4b994c.GetMemberGroupsRequestBuilder) {
    return i54c28e8ec50e0dd4eeee8073ed79cecee11ba12e3813ec711d8127af3d4b994c.NewGetMemberGroupsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetMemberObjects provides operations to call the getMemberObjects method.
func (m *UserItemRequestBuilder) GetMemberObjects()(*ide0a33943859ca09b2034a1804f60d32a108dc09bf7012d68e5d3c63ac617ddb.GetMemberObjectsRequestBuilder) {
    return ide0a33943859ca09b2034a1804f60d32a108dc09bf7012d68e5d3c63ac617ddb.NewGetMemberObjectsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// InferenceClassification provides operations to manage the inferenceClassification property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) InferenceClassification()(*i64f47462143d9c5d6613e409adca3c0fd24631de22e340968821ffdff467e657.InferenceClassificationRequestBuilder) {
    return i64f47462143d9c5d6613e409adca3c0fd24631de22e340968821ffdff467e657.NewInferenceClassificationRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Insights provides operations to manage the insights property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Insights()(*i8deae6de37c6b836f34cd5c6b8fab91a73043c26ee346f7410a8dcf88ea64c99.InsightsRequestBuilder) {
    return i8deae6de37c6b836f34cd5c6b8fab91a73043c26ee346f7410a8dcf88ea64c99.NewInsightsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// JoinedTeams provides operations to manage the joinedTeams property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) JoinedTeams()(*i7f31fc68d418ebfd772842294b63ed5b4dafc3cab499c7bf25b7634a901656bb.JoinedTeamsRequestBuilder) {
    return i7f31fc68d418ebfd772842294b63ed5b4dafc3cab499c7bf25b7634a901656bb.NewJoinedTeamsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// JoinedTeamsById provides operations to manage the joinedTeams property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) JoinedTeamsById(id string)(*i0b4dff6b82c5c98b58b9ff11aead66a7e49d0eac6b18b1490fd01bf9597cf36c.TeamItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["team%2Did"] = id
    }
    return i0b4dff6b82c5c98b58b9ff11aead66a7e49d0eac6b18b1490fd01bf9597cf36c.NewTeamItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// LicenseDetails provides operations to manage the licenseDetails property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) LicenseDetails()(*i065d0a279baf0cb315c26855c2ba7a5c0701939d8cae543a4fb9b7bb2c7cbda5.LicenseDetailsRequestBuilder) {
    return i065d0a279baf0cb315c26855c2ba7a5c0701939d8cae543a4fb9b7bb2c7cbda5.NewLicenseDetailsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// LicenseDetailsById provides operations to manage the licenseDetails property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) LicenseDetailsById(id string)(*i03f3617a0bebd4a3033a4ba021c004c146b37855bfeed1510a905b18d838b51e.LicenseDetailsItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["licenseDetails%2Did"] = id
    }
    return i03f3617a0bebd4a3033a4ba021c004c146b37855bfeed1510a905b18d838b51e.NewLicenseDetailsItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MailFolders provides operations to manage the mailFolders property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) MailFolders()(*ib946e3869f342130dddab8c76c89f264a76e899f07657349241151b4e41b4286.MailFoldersRequestBuilder) {
    return ib946e3869f342130dddab8c76c89f264a76e899f07657349241151b4e41b4286.NewMailFoldersRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MailFoldersById provides operations to manage the mailFolders property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) MailFoldersById(id string)(*id1f5db1c2220a25bede6871151a00a067652bdb75e2db4d199ca6ffcdd70b342.MailFolderItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["mailFolder%2Did"] = id
    }
    return id1f5db1c2220a25bede6871151a00a067652bdb75e2db4d199ca6ffcdd70b342.NewMailFolderItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ManagedAppRegistrations provides operations to manage the managedAppRegistrations property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) ManagedAppRegistrations()(*if12f3532feee7012e3a4dc1e3a2444ca8fab8f29e8128baafe0ef00f6ed6015c.ManagedAppRegistrationsRequestBuilder) {
    return if12f3532feee7012e3a4dc1e3a2444ca8fab8f29e8128baafe0ef00f6ed6015c.NewManagedAppRegistrationsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ManagedAppRegistrationsById provides operations to manage the managedAppRegistrations property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) ManagedAppRegistrationsById(id string)(*i5835fd9ad2a434b5df38fb1f3e5d0ec77c2f4ba2d0a321c93550f076a0281720.ManagedAppRegistrationItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["managedAppRegistration%2Did"] = id
    }
    return i5835fd9ad2a434b5df38fb1f3e5d0ec77c2f4ba2d0a321c93550f076a0281720.NewManagedAppRegistrationItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ManagedDevices provides operations to manage the managedDevices property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) ManagedDevices()(*i91af54e73e80a81d2f98c6f72fbb6c81539e67ba0480227ae8e2e5f2b8f467fa.ManagedDevicesRequestBuilder) {
    return i91af54e73e80a81d2f98c6f72fbb6c81539e67ba0480227ae8e2e5f2b8f467fa.NewManagedDevicesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ManagedDevicesById provides operations to manage the managedDevices property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) ManagedDevicesById(id string)(*i0c6b3d37dc64615f253a6fd9a58f59f536aa93a57f99bb5f9d4210ae83aa5112.ManagedDeviceItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["managedDevice%2Did"] = id
    }
    return i0c6b3d37dc64615f253a6fd9a58f59f536aa93a57f99bb5f9d4210ae83aa5112.NewManagedDeviceItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Manager provides operations to manage the manager property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Manager()(*i3fd2d5130e27e5ed481f7a0565ee17271d52e0ea15f0dec8ae436a3780d1d8d2.ManagerRequestBuilder) {
    return i3fd2d5130e27e5ed481f7a0565ee17271d52e0ea15f0dec8ae436a3780d1d8d2.NewManagerRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MemberOf provides operations to manage the memberOf property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) MemberOf()(*ia6a0f4f70d60f824df18f07f8918c6e9ac043f588293118e5d91eb921d5dffbf.MemberOfRequestBuilder) {
    return ia6a0f4f70d60f824df18f07f8918c6e9ac043f588293118e5d91eb921d5dffbf.NewMemberOfRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MemberOfById provides operations to manage the memberOf property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) MemberOfById(id string)(*i5903dc47170759db544ed9d0a0c895f70016bf360d13d02d94dd87cf4f13cd5d.DirectoryObjectItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryObject%2Did"] = id
    }
    return i5903dc47170759db544ed9d0a0c895f70016bf360d13d02d94dd87cf4f13cd5d.NewDirectoryObjectItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Messages provides operations to manage the messages property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Messages()(*i0a6d7999baeed47d1aa02960cd9f3f120528ecae4240bf52300542a8dc9f978d.MessagesRequestBuilder) {
    return i0a6d7999baeed47d1aa02960cd9f3f120528ecae4240bf52300542a8dc9f978d.NewMessagesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MessagesById provides operations to manage the messages property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) MessagesById(id string)(*if33563b345c020375260161400a6b809f00973f17d5ab870445f3c67ca7323b8.MessageItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["message%2Did"] = id
    }
    return if33563b345c020375260161400a6b809f00973f17d5ab870445f3c67ca7323b8.NewMessageItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Oauth2PermissionGrants provides operations to manage the oauth2PermissionGrants property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Oauth2PermissionGrants()(*iea339692a4ece4755a03d26004b71e7bbfe3ad28cf8036bdc4b972512eb00359.Oauth2PermissionGrantsRequestBuilder) {
    return iea339692a4ece4755a03d26004b71e7bbfe3ad28cf8036bdc4b972512eb00359.NewOauth2PermissionGrantsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Oauth2PermissionGrantsById provides operations to manage the oauth2PermissionGrants property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Oauth2PermissionGrantsById(id string)(*i49037195fa7f3333c97e06f8131490f42ed26328e1c7d9d27cba659e29c03a56.OAuth2PermissionGrantItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["oAuth2PermissionGrant%2Did"] = id
    }
    return i49037195fa7f3333c97e06f8131490f42ed26328e1c7d9d27cba659e29c03a56.NewOAuth2PermissionGrantItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Onenote provides operations to manage the onenote property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Onenote()(*ib8e8ec0d09197f1fd57167b1b67da9bcc229fa24977921ce6406f24aa89a6884.OnenoteRequestBuilder) {
    return ib8e8ec0d09197f1fd57167b1b67da9bcc229fa24977921ce6406f24aa89a6884.NewOnenoteRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OnlineMeetings provides operations to manage the onlineMeetings property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) OnlineMeetings()(*i2c0f098b49b5f2edba81f4095f0c699516263b7cd6ff316d31f469f15b8366b0.OnlineMeetingsRequestBuilder) {
    return i2c0f098b49b5f2edba81f4095f0c699516263b7cd6ff316d31f469f15b8366b0.NewOnlineMeetingsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OnlineMeetingsById provides operations to manage the onlineMeetings property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) OnlineMeetingsById(id string)(*i3ecd250b638acba95f41b3db07eaa507842a51f8e0cd3640cb9b75f0052ceff6.OnlineMeetingItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["onlineMeeting%2Did"] = id
    }
    return i3ecd250b638acba95f41b3db07eaa507842a51f8e0cd3640cb9b75f0052ceff6.NewOnlineMeetingItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Outlook provides operations to manage the outlook property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Outlook()(*i418d77dc601ca4dc2c568aef8702ce31b38463f4e1701c6c788adcbeb32ca3e8.OutlookRequestBuilder) {
    return i418d77dc601ca4dc2c568aef8702ce31b38463f4e1701c6c788adcbeb32ca3e8.NewOutlookRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OwnedDevices provides operations to manage the ownedDevices property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) OwnedDevices()(*i005f0821c1e57467ed9e31e0f65ebf8795241942e8afc0c6c9cc866dc3b70de2.OwnedDevicesRequestBuilder) {
    return i005f0821c1e57467ed9e31e0f65ebf8795241942e8afc0c6c9cc866dc3b70de2.NewOwnedDevicesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OwnedDevicesById provides operations to manage the ownedDevices property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) OwnedDevicesById(id string)(*i755269849b1d98c9e29fa525a149a4aede0223fcbc1aefdb163bc3d1fc4aa32e.DirectoryObjectItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryObject%2Did"] = id
    }
    return i755269849b1d98c9e29fa525a149a4aede0223fcbc1aefdb163bc3d1fc4aa32e.NewDirectoryObjectItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// OwnedObjects provides operations to manage the ownedObjects property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) OwnedObjects()(*i152c9659f52d4d76fff6d586e3de05874eee179419588197c21a3adc1e3a6c7b.OwnedObjectsRequestBuilder) {
    return i152c9659f52d4d76fff6d586e3de05874eee179419588197c21a3adc1e3a6c7b.NewOwnedObjectsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// OwnedObjectsById provides operations to manage the ownedObjects property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) OwnedObjectsById(id string)(*if1c7e6e4f2864e681dc58961189ef8bfe8cee67a94c24f3b9f32c1ffafdd2aa2.DirectoryObjectItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryObject%2Did"] = id
    }
    return if1c7e6e4f2864e681dc58961189ef8bfe8cee67a94c24f3b9f32c1ffafdd2aa2.NewDirectoryObjectItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update the properties of a user object. Not all properties can be updated by Member or Guest users with their default permissions without Administrator roles. Compare member and guest default permissions to see properties they can manage.
func (m *UserItemRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable, requestConfiguration *UserItemRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateUserFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Userable), nil
}
// People provides operations to manage the people property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) People()(*i2b4d48c4393811d79f4d4cc1f4addb3b49f58a53edcf7d1d7140e29a22ab7504.PeopleRequestBuilder) {
    return i2b4d48c4393811d79f4d4cc1f4addb3b49f58a53edcf7d1d7140e29a22ab7504.NewPeopleRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PeopleById provides operations to manage the people property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) PeopleById(id string)(*i8465a2d5f26666d08a0baebc755ef48f1f0dc18551f929c58c09d0110510a317.PersonItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["person%2Did"] = id
    }
    return i8465a2d5f26666d08a0baebc755ef48f1f0dc18551f929c58c09d0110510a317.NewPersonItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Photo provides operations to manage the photo property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Photo()(*icf56cc7fcccc9803d5655a975f93bee1c76d14d27e2e180c40ef7bb1c3c18a68.PhotoRequestBuilder) {
    return icf56cc7fcccc9803d5655a975f93bee1c76d14d27e2e180c40ef7bb1c3c18a68.NewPhotoRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Photos provides operations to manage the photos property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Photos()(*idea2a9c81deaaa3e33a2bcef2386387935820f4f571d71af1fc2daa24d8fa032.PhotosRequestBuilder) {
    return idea2a9c81deaaa3e33a2bcef2386387935820f4f571d71af1fc2daa24d8fa032.NewPhotosRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// PhotosById provides operations to manage the photos property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) PhotosById(id string)(*i5f2b6c8998555350079828979727c205406cf7155e907008bec42440c3e9724a.ProfilePhotoItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["profilePhoto%2Did"] = id
    }
    return i5f2b6c8998555350079828979727c205406cf7155e907008bec42440c3e9724a.NewProfilePhotoItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Planner provides operations to manage the planner property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Planner()(*i4a474d43148dd96e52e36b3ed7c52a7d75f46d95adf4fa0ce96ace2509d33744.PlannerRequestBuilder) {
    return i4a474d43148dd96e52e36b3ed7c52a7d75f46d95adf4fa0ce96ace2509d33744.NewPlannerRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Presence provides operations to manage the presence property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Presence()(*i0a4c7a84985f2276fcfbcf69e2fe37b80ef87c2622cc8763371022616974902d.PresenceRequestBuilder) {
    return i0a4c7a84985f2276fcfbcf69e2fe37b80ef87c2622cc8763371022616974902d.NewPresenceRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RegisteredDevices provides operations to manage the registeredDevices property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) RegisteredDevices()(*i631d94a7300d4142a2d7a4c4d0d424e98d3fb058f106e6567a8ac6476f5956b6.RegisteredDevicesRequestBuilder) {
    return i631d94a7300d4142a2d7a4c4d0d424e98d3fb058f106e6567a8ac6476f5956b6.NewRegisteredDevicesRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RegisteredDevicesById provides operations to manage the registeredDevices property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) RegisteredDevicesById(id string)(*i3468dc539adf033fc51d31c6b79eb77caea5661d9461c49781e068c23230b9fc.DirectoryObjectItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryObject%2Did"] = id
    }
    return i3468dc539adf033fc51d31c6b79eb77caea5661d9461c49781e068c23230b9fc.NewDirectoryObjectItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// ReminderViewWithStartDateTimeWithEndDateTime provides operations to call the reminderView method.
func (m *UserItemRequestBuilder) ReminderViewWithStartDateTimeWithEndDateTime(endDateTime *string, startDateTime *string)(*iddc156b6ba287ee33d5592797e8da64087fd068cf179bd892944253be67ac174.ReminderViewWithStartDateTimeWithEndDateTimeRequestBuilder) {
    return iddc156b6ba287ee33d5592797e8da64087fd068cf179bd892944253be67ac174.NewReminderViewWithStartDateTimeWithEndDateTimeRequestBuilderInternal(m.pathParameters, m.requestAdapter, endDateTime, startDateTime);
}
// RemoveAllDevicesFromManagement provides operations to call the removeAllDevicesFromManagement method.
func (m *UserItemRequestBuilder) RemoveAllDevicesFromManagement()(*i6d278e5cb359938f2362c090cc322d0c80867015903a39849ea478a6cda51689.RemoveAllDevicesFromManagementRequestBuilder) {
    return i6d278e5cb359938f2362c090cc322d0c80867015903a39849ea478a6cda51689.NewRemoveAllDevicesFromManagementRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ReprocessLicenseAssignment provides operations to call the reprocessLicenseAssignment method.
func (m *UserItemRequestBuilder) ReprocessLicenseAssignment()(*ifd940c65ea13f1270f41f753ee5e7761f074915bd6a6c64b2adb77f70aea8c32.ReprocessLicenseAssignmentRequestBuilder) {
    return ifd940c65ea13f1270f41f753ee5e7761f074915bd6a6c64b2adb77f70aea8c32.NewReprocessLicenseAssignmentRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Restore provides operations to call the restore method.
func (m *UserItemRequestBuilder) Restore()(*i9f87171771d54744395d1043beb42c14430d7cb4f1557826ebef7c5b2584f7cd.RestoreRequestBuilder) {
    return i9f87171771d54744395d1043beb42c14430d7cb4f1557826ebef7c5b2584f7cd.NewRestoreRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// RevokeSignInSessions provides operations to call the revokeSignInSessions method.
func (m *UserItemRequestBuilder) RevokeSignInSessions()(*i1090e84ad9b76c4c279467ebbc90ff155a540b45b098e680706d2aae300639eb.RevokeSignInSessionsRequestBuilder) {
    return i1090e84ad9b76c4c279467ebbc90ff155a540b45b098e680706d2aae300639eb.NewRevokeSignInSessionsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ScopedRoleMemberOf provides operations to manage the scopedRoleMemberOf property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) ScopedRoleMemberOf()(*ia99f649e03b3020696f5c2c11f42a9ad7f094cb80add889863fc7947d3bf9ec7.ScopedRoleMemberOfRequestBuilder) {
    return ia99f649e03b3020696f5c2c11f42a9ad7f094cb80add889863fc7947d3bf9ec7.NewScopedRoleMemberOfRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ScopedRoleMemberOfById provides operations to manage the scopedRoleMemberOf property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) ScopedRoleMemberOfById(id string)(*i983bdc3dce2d947aec82e1e9dd5edd887e8d861d589287d050f96af67abb1e1f.ScopedRoleMembershipItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["scopedRoleMembership%2Did"] = id
    }
    return i983bdc3dce2d947aec82e1e9dd5edd887e8d861d589287d050f96af67abb1e1f.NewScopedRoleMembershipItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// SendMail provides operations to call the sendMail method.
func (m *UserItemRequestBuilder) SendMail()(*i16f13fb20dc21e9cbd93cbca432ac5cb871c47a97e5bebce67f5eab6fe8276e9.SendMailRequestBuilder) {
    return i16f13fb20dc21e9cbd93cbca432ac5cb871c47a97e5bebce67f5eab6fe8276e9.NewSendMailRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Settings provides operations to manage the settings property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Settings()(*i225d659ae554743be8acfa963ab6384d1f6e82911866fa2e3c01918b6a0637b1.SettingsRequestBuilder) {
    return i225d659ae554743be8acfa963ab6384d1f6e82911866fa2e3c01918b6a0637b1.NewSettingsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Teamwork provides operations to manage the teamwork property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Teamwork()(*i94ff5ee566023b18b8b2ac8c1da92081f3bf4520205017599c00aed9faff56de.TeamworkRequestBuilder) {
    return i94ff5ee566023b18b8b2ac8c1da92081f3bf4520205017599c00aed9faff56de.NewTeamworkRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Todo provides operations to manage the todo property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) Todo()(*i61d84fd6521641b94f790414a5b9abccf870ea66012c270c06f3ea7df770aad8.TodoRequestBuilder) {
    return i61d84fd6521641b94f790414a5b9abccf870ea66012c270c06f3ea7df770aad8.NewTodoRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TransitiveMemberOf provides operations to manage the transitiveMemberOf property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) TransitiveMemberOf()(*i11db6285564609d0b36b9032e0d8f6d257828b887e558375da4354121f40a3d0.TransitiveMemberOfRequestBuilder) {
    return i11db6285564609d0b36b9032e0d8f6d257828b887e558375da4354121f40a3d0.NewTransitiveMemberOfRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// TransitiveMemberOfById provides operations to manage the transitiveMemberOf property of the microsoft.graph.user entity.
func (m *UserItemRequestBuilder) TransitiveMemberOfById(id string)(*ie9198f71b9fb192ecc84221b88715d18fc05b44b87d4b863d0e98c3c6b88c246.DirectoryObjectItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["directoryObject%2Did"] = id
    }
    return ie9198f71b9fb192ecc84221b88715d18fc05b44b87d4b863d0e98c3c6b88c246.NewDirectoryObjectItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// TranslateExchangeIds provides operations to call the translateExchangeIds method.
func (m *UserItemRequestBuilder) TranslateExchangeIds()(*ifdf31129bec8038d82607a6d973830e0b090629373d5812152e5046514a4bfb0.TranslateExchangeIdsRequestBuilder) {
    return ifdf31129bec8038d82607a6d973830e0b090629373d5812152e5046514a4bfb0.NewTranslateExchangeIdsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// WipeManagedAppRegistrationsByDeviceTag provides operations to call the wipeManagedAppRegistrationsByDeviceTag method.
func (m *UserItemRequestBuilder) WipeManagedAppRegistrationsByDeviceTag()(*i4af8697d288ced3f62e5f7467084c929a355a7f4c2140a3f51072b00f73663c1.WipeManagedAppRegistrationsByDeviceTagRequestBuilder) {
    return i4af8697d288ced3f62e5f7467084c929a355a7f4c2140a3f51072b00f73663c1.NewWipeManagedAppRegistrationsByDeviceTagRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
