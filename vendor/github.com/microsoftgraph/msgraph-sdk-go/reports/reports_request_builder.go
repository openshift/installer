package reports

import (
    "context"
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
    i01986228cfa34d663610609938a20fe1aa899eb3ae22c44d689bf3a7cc4fcec2 "github.com/microsoftgraph/msgraph-sdk-go/reports/getonedriveusageaccountdetailwithdate"
    i028952a5521bf457c8a8796489c7f5a1137cf7b0745234478ff5fd12bb087564 "github.com/microsoftgraph/msgraph-sdk-go/reports/getsharepointsiteusagestoragewithperiod"
    i037933f7467ae74beac89af9fadc6cf1a5a8272d8f21a5c2f98bc8f8e7c28ab6 "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinessdeviceusageuserdetailwithdate"
    i03d33a99775bdfdb348179d972066b208c7410da4d93f4776c48bb0c7619f539 "github.com/microsoftgraph/msgraph-sdk-go/reports/getyammerdeviceusagedistributionusercountswithperiod"
    i03fae197628d424772396da17b5102e4676ffef8400787778663d2f1e8b62054 "github.com/microsoftgraph/msgraph-sdk-go/reports/getemailappusageusercountswithperiod"
    i05c8d276b86beee65c00c115937a46c2753c7d328e142fe9afab6ef90b5a1a79 "github.com/microsoftgraph/msgraph-sdk-go/reports/getteamsdeviceusagedistributionusercountswithperiod"
    i05ca24f1e64f29c556c3c67a9024eea8b3aac4d350f05ff2c52f4511bf0b6f33 "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinessparticipantactivityusercountswithperiod"
    i066e5ea9f7e2ed1b047023128aab6caef0db934d50a71df44b7b6cf7bfb5a10e "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinessorganizeractivityusercountswithperiod"
    i069326085999a3926845312c56f68613c15a7dcb468bf7f0592670ea8f5896aa "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinessactivityusercountswithperiod"
    i0c4b2ea99f0789f3909bf8d293ade8bd9216c7010a8e0dc72ceea54c21c1d481 "github.com/microsoftgraph/msgraph-sdk-go/reports/getteamsuseractivityuserdetailwithdate"
    i0d6612898baf60cf3ad09e5ec9209afdfaffc25ddef2fe3fddf028f8f356d2e1 "github.com/microsoftgraph/msgraph-sdk-go/reports/getoffice365activeuserdetailwithdate"
    i0e6a5a79552b1c4826c9418d405e7420746bbd08e606bfd5c84ba1ad674aed8f "github.com/microsoftgraph/msgraph-sdk-go/reports/manageddeviceenrollmentfailuredetails"
    i13fef70f90226be585263be7022fe29fa384901ba2188c080017c9272ca23e99 "github.com/microsoftgraph/msgraph-sdk-go/reports/getyammeractivityusercountswithperiod"
    i1bfeeda14fdcd3c4bc1c9b02a4d925837c77654498184deaffbfeb8f496c10c5 "github.com/microsoftgraph/msgraph-sdk-go/reports/getsharepointactivityfilecountswithperiod"
    i1ccdf78d8963c250d79813dd0917eaaaee794135fc39e35801b01ea5f2005034 "github.com/microsoftgraph/msgraph-sdk-go/reports/getonedriveactivityusercountswithperiod"
    i1d4731cdb4991f39251caf5d5cbcc8c68036c8b25a8c59c7ee8b24f06ba15ead "github.com/microsoftgraph/msgraph-sdk-go/reports/getm365appuserdetailwithperiod"
    i1d9fef6f0fe1efedede6850168705608b3b5be81427b11c9c3d64b7dc8c126ed "github.com/microsoftgraph/msgraph-sdk-go/reports/dailyprintusagebyuser"
    i200c880280f3a34037634650d2cd9969f9a07037f773af21ebb6141e40dde23b "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinesspeertopeeractivitycountswithperiod"
    i21ef6a6d7bc99608d2279d5cc2343f5300ae76d3d739202a0c1bc4ff2d13a391 "github.com/microsoftgraph/msgraph-sdk-go/reports/getoffice365groupsactivityfilecountswithperiod"
    i2306eba962c786175114d00783028649abfdfb48f769241b2a7dca15db108361 "github.com/microsoftgraph/msgraph-sdk-go/reports/getoffice365servicesusercountswithperiod"
    i251d35954b406eaf53009c768700ddbf5c49cf2961b8eba3836e2905c19fe73a "github.com/microsoftgraph/msgraph-sdk-go/reports/deviceconfigurationuseractivity"
    i2699e05c2942728d3c9dbbec2f7da6269f18ec269bdc96532cb49ca45ea2e03f "github.com/microsoftgraph/msgraph-sdk-go/reports/getsharepointsiteusagefilecountswithperiod"
    i26d1169a1f986183917395a64e40bb9c1870a54b457336ef64e7efa21c214a15 "github.com/microsoftgraph/msgraph-sdk-go/reports/getoffice365groupsactivitystoragewithperiod"
    i285bf21395498d7467b7e1d3f62e43bbbab40d875429e2800566739471aaa586 "github.com/microsoftgraph/msgraph-sdk-go/reports/getemailappusageuserdetailwithperiod"
    i2e849a1f15987664ba1a6648f0c1d0eb9e2ca5965a632cc3761ad14bb91674bf "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinessorganizeractivitycountswithperiod"
    i2ee11973211f7fc9312ca00d777ef8304aada9d21bbf7c226dda5bef84678886 "github.com/microsoftgraph/msgraph-sdk-go/reports/security"
    i2f04ee022b2f4dd0c31e5c7f89cc28c7d1e36585beb573c9c689634579fb8999 "github.com/microsoftgraph/msgraph-sdk-go/reports/getsharepointsiteusagepageswithperiod"
    i316dae3dc0a83b965993ba7f5a9d3b29d2b5dcbd706e8977d4a6cb34548c5275 "github.com/microsoftgraph/msgraph-sdk-go/reports/getonedriveusageaccountdetailwithperiod"
    i347d730506a06c407296d7557c653293a35a9c38b50c7256c1c7abb4c736ef9a "github.com/microsoftgraph/msgraph-sdk-go/reports/getsharepointsiteusagesitecountswithperiod"
    i34861f91c890c43257a1e1386ba24de1ffe44f842eadf95d13de30d27a852991 "github.com/microsoftgraph/msgraph-sdk-go/reports/getteamsuseractivityusercountswithperiod"
    i374eb5f5ea7befd0f519a8a1c2886282c8ae9910332bd6b507ad8b29eddfdf2d "github.com/microsoftgraph/msgraph-sdk-go/reports/manageddeviceenrollmentfailuredetailswithskipwithtopwithfilterwithskiptoken"
    i37cb418ef8dbdf59e2a836a306ab7fe4b5933e6211360de8f4f313f68496d242 "github.com/microsoftgraph/msgraph-sdk-go/reports/getemailactivityuserdetailwithdate"
    i37d455e922500fc7c5e44df5cc84a103ddbb8470a53860cb0c02e64fd6e8e775 "github.com/microsoftgraph/msgraph-sdk-go/reports/getmailboxusagequotastatusmailboxcountswithperiod"
    i37efe63b21f1e137ba9fdf350d68a656b504d390fe52651da7739ce398ec31c8 "github.com/microsoftgraph/msgraph-sdk-go/reports/getoffice365activationcounts"
    i39400581189e192e8c3927b2ac378c799b30d5cb2c7a932a3e2f9802d2a2f95d "github.com/microsoftgraph/msgraph-sdk-go/reports/getgrouparchivedprintjobswithgroupidwithstartdatetimewithenddatetime"
    i39d4923038966e0aa0bf413db480c947e66bb5fa71c83d893fbfb4aaacc6cfe9 "github.com/microsoftgraph/msgraph-sdk-go/reports/deviceconfigurationdeviceactivity"
    i3ed9439789e0857666ebca3cd16359b4745ffe979a682ea73e6cd26f8054ae21 "github.com/microsoftgraph/msgraph-sdk-go/reports/getprinterarchivedprintjobswithprinteridwithstartdatetimewithenddatetime"
    i3f12902a7f7457a60f97014c6db3b2bd89dac8cabe4d467b766677db3d6029d0 "github.com/microsoftgraph/msgraph-sdk-go/reports/monthlyprintusagebyuser"
    i3f4aebd4163df4cb0a4f38097f15a9f8ba73853688cf5a50b93a222eb4832abe "github.com/microsoftgraph/msgraph-sdk-go/reports/manageddeviceenrollmenttopfailureswithperiod"
    i40f65a5b54f01dfaea303a3582e82bcdaa71c2195cc00b7915431e77f38e9a0f "github.com/microsoftgraph/msgraph-sdk-go/reports/getsharepointactivityusercountswithperiod"
    i416e8a1993d00c5e5b47ab655f48ac49786eb1442a5459b7cb7023ac34c28093 "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinesspeertopeeractivityusercountswithperiod"
    i4319947bf5469fb8e87f0bc08088c383a13726c5be51ca3558867ea6abecc315 "github.com/microsoftgraph/msgraph-sdk-go/reports/getoffice365groupsactivitydetailwithdate"
    i4ea09032e0c277ad5c699279a1e67097201d217b8d0aac6b5962420831d787bd "github.com/microsoftgraph/msgraph-sdk-go/reports/getoffice365activeusercountswithperiod"
    i4ea8142505b6b02804cdd135578592c59e1ebd8c244878c8a27bac9df783863b "github.com/microsoftgraph/msgraph-sdk-go/reports/getemailappusageversionsusercountswithperiod"
    i5277312ad2cbf60385fe9448b03ae55d7aefafa14a9c09b135b745563e0dd271 "github.com/microsoftgraph/msgraph-sdk-go/reports/getyammergroupsactivitydetailwithperiod"
    i5445cb98613918d3a7a579fddf521b041c55f06714612e2097ed886a37be4286 "github.com/microsoftgraph/msgraph-sdk-go/reports/getoffice365groupsactivitygroupcountswithperiod"
    i54533b65ae17487eb09581591352cc93d62970fd3000b0af7a1abc12bf038091 "github.com/microsoftgraph/msgraph-sdk-go/reports/getemailactivityusercountswithperiod"
    i55b749bfea07392ccbd597dc8d76a153720e39963b1e66527cac5bb79f8c6b29 "github.com/microsoftgraph/msgraph-sdk-go/reports/getyammergroupsactivitygroupcountswithperiod"
    i56c28dc52ab4ecbfe37baebbfaf6f279f8d130872c0931a194264028d89e2c96 "github.com/microsoftgraph/msgraph-sdk-go/reports/getsharepointactivityuserdetailwithperiod"
    i5be64ff627afac918209c9e795b4c7ad94cf40584d7376fd0bddf4902f737286 "github.com/microsoftgraph/msgraph-sdk-go/reports/getoffice365activationsusercounts"
    i61a78ce1e0b6e2e60e6d50fe39c731a201868c12c59871387bf62c67fa03cc24 "github.com/microsoftgraph/msgraph-sdk-go/reports/getmailboxusagedetailwithperiod"
    i677bcf76c9544b6088f18b0d84351cd554f0537738c43ca67e647fffec66978c "github.com/microsoftgraph/msgraph-sdk-go/reports/getteamsdeviceusageuserdetailwithdate"
    i67d417cc0a32f095a4d9872596c73a5670bab5ca99f2c39b690a3dc3e53dacb6 "github.com/microsoftgraph/msgraph-sdk-go/reports/getoffice365activationsuserdetail"
    i6a26bcea122cbd3594cacac16b0715042c3ca45bb5f889d1fd0bd83d243a6fe3 "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinessactivityuserdetailwithdate"
    i6b3e5d0662b4c04bf1f969a3c62a48d50c4f7ecb57a3913782ce351d87a0035e "github.com/microsoftgraph/msgraph-sdk-go/reports/getsharepointactivitypageswithperiod"
    i6c58964b36e86bf47df5ddf5c8161ca0e4621aee61aaa9517f1dded3120cfc8d "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinessdeviceusageuserdetailwithperiod"
    i6e07bd3dadcc8aa2d52f07e77f1f35672ac7ce62aa2917c7a9e4e77b8f4c4db5 "github.com/microsoftgraph/msgraph-sdk-go/reports/getyammerdeviceusageuserdetailwithperiod"
    i6e8f025e6b5802865cf6c583d64d93f7e1c603ead045ca53bef1f6b9cde8a5b0 "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinessparticipantactivitycountswithperiod"
    i6edf15eaf951cb8d179890d37ee652592111a5bfeb9fd5709d3609c7593bed81 "github.com/microsoftgraph/msgraph-sdk-go/reports/getoffice365activeuserdetailwithperiod"
    i700e0ab2f1ea334038b3293b2396e12432ea3545b07ed127bb12d122bbb6bf90 "github.com/microsoftgraph/msgraph-sdk-go/reports/getmailboxusagemailboxcountswithperiod"
    i753582442e8209033f45d1a367493028141a79145904e5dbd0343dc59d62d8c9 "github.com/microsoftgraph/msgraph-sdk-go/reports/getm365appuserdetailwithdate"
    i75a65b7f6cf316241684f8aeba4d08fb2c2e04274ef9390b6003ff2a3fe87418 "github.com/microsoftgraph/msgraph-sdk-go/reports/getonedriveusageaccountcountswithperiod"
    i7774a89bdd2b6bc2f70c1e7c58f236fd1f3aa68a40a72db4898fd9c5d4f66988 "github.com/microsoftgraph/msgraph-sdk-go/reports/getoffice365groupsactivitydetailwithperiod"
    i7c11427d5ee6b6eb1e984b209bd6cb9cd7704d1573048256590bf48fe10c95c2 "github.com/microsoftgraph/msgraph-sdk-go/reports/getemailappusageappsusercountswithperiod"
    i7ee204dcbd579a0c0d252ddbfcdd01bf9f244b2b38f671f02c4108f867a99c73 "github.com/microsoftgraph/msgraph-sdk-go/reports/getyammergroupsactivitycountswithperiod"
    i8ad77c0341e431729102d8b93f9066b82a7f4fc8ee3e748eeea007bf3d9d5b18 "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinessactivitycountswithperiod"
    i8d02e7a0b578c5b4465f96608d724ae82f61e67266f11ff2ab8aad96d55d1ac1 "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinessorganizeractivityminutecountswithperiod"
    i9b65430927e689a26a700312296e423fc6e6e6c0bb86f72287b75c7915e695e0 "github.com/microsoftgraph/msgraph-sdk-go/reports/getteamsuseractivitycountswithperiod"
    i9b77e3c3d15971fe31fa9e23d6bdc4c74a584c922585117306c4db805ac4131c "github.com/microsoftgraph/msgraph-sdk-go/reports/getyammerdeviceusageusercountswithperiod"
    i9d1992f642df096181dd0b243f5639c29682f838689705077b65251f267d222f "github.com/microsoftgraph/msgraph-sdk-go/reports/manageddeviceenrollmenttopfailures"
    i9de2dd5236ac703b10206d8522e5a8f5490288eadb62a3fbf88d9becd9f061bf "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinessdeviceusageusercountswithperiod"
    ia3da511007b4da9fa7b31d42bcef05972fa35df25142f53ebb0f2944e3932886 "github.com/microsoftgraph/msgraph-sdk-go/reports/getyammergroupsactivitydetailwithdate"
    ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a "github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
    ia5f41691c1fbeb62b6b1de131fa27f7306db515c8759a9e49ec88553c0e6d97d "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinessdeviceusagedistributionusercountswithperiod"
    ia79ac9c43dc0218a9c3f469a247120979433dd0e3e53de3949570224b3cf3936 "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinessactivityuserdetailwithperiod"
    iaaacf2b7aacf4d5f8107fa18d6f0d6cb33066e3335d8ab8cc1a2adfa1f4e3c42 "github.com/microsoftgraph/msgraph-sdk-go/reports/getteamsdeviceusageuserdetailwithperiod"
    iab84fe97b770990c561344b8be4caaca7e6d962145fd1ad9fc9e851313b38857 "github.com/microsoftgraph/msgraph-sdk-go/reports/getonedriveactivityuserdetailwithperiod"
    iafadefc9c12601a3611573bb976b73b251f9ea476a68e3bdc9ceab9e0d5cc308 "github.com/microsoftgraph/msgraph-sdk-go/reports/getuserarchivedprintjobswithuseridwithstartdatetimewithenddatetime"
    ib3bef6d1d1d843d8e20e58b41c29e9b82bb511ac97e381039f942ffee79a3b9a "github.com/microsoftgraph/msgraph-sdk-go/reports/getsharepointactivityuserdetailwithdate"
    ib9c6325caeea63e2dd57f99af77223cbb370ce799b27f8b8a3f6a3f987a9f9fe "github.com/microsoftgraph/msgraph-sdk-go/reports/getsharepointsiteusagedetailwithperiod"
    ibc84667b2a2b840609232b734d0a2873e66315fbe63aeadb3c3d468ad9f84d75 "github.com/microsoftgraph/msgraph-sdk-go/reports/getonedriveusagefilecountswithperiod"
    ibfd05176830665cf1dae7366e8b491572f134d04e12120a5499c283d6c8ae062 "github.com/microsoftgraph/msgraph-sdk-go/reports/getteamsuseractivityuserdetailwithperiod"
    ic5508dd2223dfecdf5814372064e28b262f8f8bfe997c49795b7f76e30889837 "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinesspeertopeeractivityminutecountswithperiod"
    ic5ac4372354eec289834ac2ae20a4b14244024792aa3c3ad01b0f0a204b1e81b "github.com/microsoftgraph/msgraph-sdk-go/reports/getsharepointsiteusagedetailwithdate"
    ic78e2d968f685321dafd3ca273401e484836d7ca84e4da07ae045edca4ca55a2 "github.com/microsoftgraph/msgraph-sdk-go/reports/dailyprintusagebyprinter"
    id05b47283d832e55b732685b30d918ee0e303ae639c6447cfff5c7a7089741c7 "github.com/microsoftgraph/msgraph-sdk-go/reports/getmailboxusagestoragewithperiod"
    id379a8afbc53e3411a43fbddcd71ea21f90a2a71ada609569a5c9040ebee5f51 "github.com/microsoftgraph/msgraph-sdk-go/reports/getyammeractivityuserdetailwithperiod"
    id7d2a94ac1ea97cc24d2cd9cb47182b168535107d0d09ad275c4e955f417e029 "github.com/microsoftgraph/msgraph-sdk-go/reports/getyammeractivityuserdetailwithdate"
    idcca069eb97aac5c90ca4ddb83fee2301a14619967cfdadbb5d0a81effc5cc59 "github.com/microsoftgraph/msgraph-sdk-go/reports/getskypeforbusinessparticipantactivityminutecountswithperiod"
    ie49c06630ffed10bc5667592819d2e6ba25594d178187e0c69a98a6cc5ca2e88 "github.com/microsoftgraph/msgraph-sdk-go/reports/getemailappusageuserdetailwithdate"
    ie603b5f7e630bd6d2e83589da6d79d43b190e1d4c7fc28656991ffe48fca4828 "github.com/microsoftgraph/msgraph-sdk-go/reports/getm365appplatformusercountswithperiod"
    ie6b56eee43945ce24ada2464a67b7217385822eb4143ed753098debeb7fa82eb "github.com/microsoftgraph/msgraph-sdk-go/reports/getemailactivityuserdetailwithperiod"
    iec7256d00a8e37c8dec628d1626f62f1224b45c25e3c24418830b8fede7bc006 "github.com/microsoftgraph/msgraph-sdk-go/reports/getonedriveactivityfilecountswithperiod"
    ieef21740e587e5f9d695c0f6753ee2fd4da00a0ced4e73355cae511e81a4e3cf "github.com/microsoftgraph/msgraph-sdk-go/reports/getemailactivitycountswithperiod"
    if1205fb664936127831759fefe4c2c43852474a3fc0419f51e89975f0f281d32 "github.com/microsoftgraph/msgraph-sdk-go/reports/monthlyprintusagebyprinter"
    if97af71691601e0f0046051c8e82fb71f8db2f1d3b5f8e7e40346d3fa47b196f "github.com/microsoftgraph/msgraph-sdk-go/reports/getonedriveusagestoragewithperiod"
    if97f8cb01f2d6a02697169f75d571d9a67fc56f2041590c23f83bf2114377bda "github.com/microsoftgraph/msgraph-sdk-go/reports/getyammerdeviceusageuserdetailwithdate"
    ifb70744a42eac9e3054301d8328a61f22ab480c9dcaaaa3f398cf2a26e1b242b "github.com/microsoftgraph/msgraph-sdk-go/reports/getoffice365groupsactivitycountswithperiod"
    ifd36d48788fa091c2ebd9158094174f0f9a871398672c26de206aab3755cc266 "github.com/microsoftgraph/msgraph-sdk-go/reports/getm365appusercountswithperiod"
    ifd45221e0ff465bb6f37fa740578c697634181cf333e3c01e4db49019e2ba5eb "github.com/microsoftgraph/msgraph-sdk-go/reports/getyammeractivitycountswithperiod"
    ife6d24553cc77853635c06a75c3587d7c978531be356ad656116b7742f60e1f7 "github.com/microsoftgraph/msgraph-sdk-go/reports/getonedriveactivityuserdetailwithdate"
    ife99981e20f3609c2111e3e2aff666b997075ba1e7c98ff44f1950f1d6265ccb "github.com/microsoftgraph/msgraph-sdk-go/reports/getteamsdeviceusageusercountswithperiod"
    i0fbd1879eec3f09756dabd53c9de402f839fcf5ed7e8be0da0d3e4135b124b0a "github.com/microsoftgraph/msgraph-sdk-go/reports/monthlyprintusagebyprinter/item"
    i48adffda977c4642ea9d1490fef3ceaa45cfab2defb01a97672b497f12eee0bc "github.com/microsoftgraph/msgraph-sdk-go/reports/dailyprintusagebyprinter/item"
    i80d61621ea489f37484a66bdee6bdaa9bbb8a0b83a1cc326a225bedfc4c6dd91 "github.com/microsoftgraph/msgraph-sdk-go/reports/dailyprintusagebyuser/item"
    if4346a9921b148666613b317b233f49301e4eef2f130f6b9e02414250e5a2294 "github.com/microsoftgraph/msgraph-sdk-go/reports/monthlyprintusagebyuser/item"
)

// ReportsRequestBuilder provides operations to manage the reportRoot singleton.
type ReportsRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ReportsRequestBuilderGetQueryParameters get reports
type ReportsRequestBuilderGetQueryParameters struct {
    // Expand related entities
    Expand []string `uriparametername:"%24expand"`
    // Select properties to be returned
    Select []string `uriparametername:"%24select"`
}
// ReportsRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ReportsRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *ReportsRequestBuilderGetQueryParameters
}
// ReportsRequestBuilderPatchRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ReportsRequestBuilderPatchRequestConfiguration struct {
    // Request headers
    Headers map[string]string
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewReportsRequestBuilderInternal instantiates a new ReportsRequestBuilder and sets the default values.
func NewReportsRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ReportsRequestBuilder) {
    m := &ReportsRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/reports{?%24select,%24expand}";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewReportsRequestBuilder instantiates a new ReportsRequestBuilder and sets the default values.
func NewReportsRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ReportsRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewReportsRequestBuilderInternal(urlParams, requestAdapter)
}
// CreateGetRequestInformation get reports
func (m *ReportsRequestBuilder) CreateGetRequestInformation(ctx context.Context, requestConfiguration *ReportsRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// CreatePatchRequestInformation update reports
func (m *ReportsRequestBuilder) CreatePatchRequestInformation(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ReportRootable, requestConfiguration *ReportsRequestBuilderPatchRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
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
// DailyPrintUsageByPrinter provides operations to manage the dailyPrintUsageByPrinter property of the microsoft.graph.reportRoot entity.
func (m *ReportsRequestBuilder) DailyPrintUsageByPrinter()(*ic78e2d968f685321dafd3ca273401e484836d7ca84e4da07ae045edca4ca55a2.DailyPrintUsageByPrinterRequestBuilder) {
    return ic78e2d968f685321dafd3ca273401e484836d7ca84e4da07ae045edca4ca55a2.NewDailyPrintUsageByPrinterRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DailyPrintUsageByPrinterById provides operations to manage the dailyPrintUsageByPrinter property of the microsoft.graph.reportRoot entity.
func (m *ReportsRequestBuilder) DailyPrintUsageByPrinterById(id string)(*i48adffda977c4642ea9d1490fef3ceaa45cfab2defb01a97672b497f12eee0bc.PrintUsageByPrinterItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["printUsageByPrinter%2Did"] = id
    }
    return i48adffda977c4642ea9d1490fef3ceaa45cfab2defb01a97672b497f12eee0bc.NewPrintUsageByPrinterItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DailyPrintUsageByUser provides operations to manage the dailyPrintUsageByUser property of the microsoft.graph.reportRoot entity.
func (m *ReportsRequestBuilder) DailyPrintUsageByUser()(*i1d9fef6f0fe1efedede6850168705608b3b5be81427b11c9c3d64b7dc8c126ed.DailyPrintUsageByUserRequestBuilder) {
    return i1d9fef6f0fe1efedede6850168705608b3b5be81427b11c9c3d64b7dc8c126ed.NewDailyPrintUsageByUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DailyPrintUsageByUserById provides operations to manage the dailyPrintUsageByUser property of the microsoft.graph.reportRoot entity.
func (m *ReportsRequestBuilder) DailyPrintUsageByUserById(id string)(*i80d61621ea489f37484a66bdee6bdaa9bbb8a0b83a1cc326a225bedfc4c6dd91.PrintUsageByUserItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["printUsageByUser%2Did"] = id
    }
    return i80d61621ea489f37484a66bdee6bdaa9bbb8a0b83a1cc326a225bedfc4c6dd91.NewPrintUsageByUserItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// DeviceConfigurationDeviceActivity provides operations to call the deviceConfigurationDeviceActivity method.
func (m *ReportsRequestBuilder) DeviceConfigurationDeviceActivity()(*i39d4923038966e0aa0bf413db480c947e66bb5fa71c83d893fbfb4aaacc6cfe9.DeviceConfigurationDeviceActivityRequestBuilder) {
    return i39d4923038966e0aa0bf413db480c947e66bb5fa71c83d893fbfb4aaacc6cfe9.NewDeviceConfigurationDeviceActivityRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// DeviceConfigurationUserActivity provides operations to call the deviceConfigurationUserActivity method.
func (m *ReportsRequestBuilder) DeviceConfigurationUserActivity()(*i251d35954b406eaf53009c768700ddbf5c49cf2961b8eba3836e2905c19fe73a.DeviceConfigurationUserActivityRequestBuilder) {
    return i251d35954b406eaf53009c768700ddbf5c49cf2961b8eba3836e2905c19fe73a.NewDeviceConfigurationUserActivityRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// Get get reports
func (m *ReportsRequestBuilder) Get(ctx context.Context, requestConfiguration *ReportsRequestBuilderGetRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ReportRootable, error) {
    requestInfo, err := m.CreateGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateReportRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ReportRootable), nil
}
// GetEmailActivityCountsWithPeriod provides operations to call the getEmailActivityCounts method.
func (m *ReportsRequestBuilder) GetEmailActivityCountsWithPeriod(period *string)(*ieef21740e587e5f9d695c0f6753ee2fd4da00a0ced4e73355cae511e81a4e3cf.GetEmailActivityCountsWithPeriodRequestBuilder) {
    return ieef21740e587e5f9d695c0f6753ee2fd4da00a0ced4e73355cae511e81a4e3cf.NewGetEmailActivityCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetEmailActivityUserCountsWithPeriod provides operations to call the getEmailActivityUserCounts method.
func (m *ReportsRequestBuilder) GetEmailActivityUserCountsWithPeriod(period *string)(*i54533b65ae17487eb09581591352cc93d62970fd3000b0af7a1abc12bf038091.GetEmailActivityUserCountsWithPeriodRequestBuilder) {
    return i54533b65ae17487eb09581591352cc93d62970fd3000b0af7a1abc12bf038091.NewGetEmailActivityUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetEmailActivityUserDetailWithDate provides operations to call the getEmailActivityUserDetail method.
func (m *ReportsRequestBuilder) GetEmailActivityUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*i37cb418ef8dbdf59e2a836a306ab7fe4b5933e6211360de8f4f313f68496d242.GetEmailActivityUserDetailWithDateRequestBuilder) {
    return i37cb418ef8dbdf59e2a836a306ab7fe4b5933e6211360de8f4f313f68496d242.NewGetEmailActivityUserDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetEmailActivityUserDetailWithPeriod provides operations to call the getEmailActivityUserDetail method.
func (m *ReportsRequestBuilder) GetEmailActivityUserDetailWithPeriod(period *string)(*ie6b56eee43945ce24ada2464a67b7217385822eb4143ed753098debeb7fa82eb.GetEmailActivityUserDetailWithPeriodRequestBuilder) {
    return ie6b56eee43945ce24ada2464a67b7217385822eb4143ed753098debeb7fa82eb.NewGetEmailActivityUserDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetEmailAppUsageAppsUserCountsWithPeriod provides operations to call the getEmailAppUsageAppsUserCounts method.
func (m *ReportsRequestBuilder) GetEmailAppUsageAppsUserCountsWithPeriod(period *string)(*i7c11427d5ee6b6eb1e984b209bd6cb9cd7704d1573048256590bf48fe10c95c2.GetEmailAppUsageAppsUserCountsWithPeriodRequestBuilder) {
    return i7c11427d5ee6b6eb1e984b209bd6cb9cd7704d1573048256590bf48fe10c95c2.NewGetEmailAppUsageAppsUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetEmailAppUsageUserCountsWithPeriod provides operations to call the getEmailAppUsageUserCounts method.
func (m *ReportsRequestBuilder) GetEmailAppUsageUserCountsWithPeriod(period *string)(*i03fae197628d424772396da17b5102e4676ffef8400787778663d2f1e8b62054.GetEmailAppUsageUserCountsWithPeriodRequestBuilder) {
    return i03fae197628d424772396da17b5102e4676ffef8400787778663d2f1e8b62054.NewGetEmailAppUsageUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetEmailAppUsageUserDetailWithDate provides operations to call the getEmailAppUsageUserDetail method.
func (m *ReportsRequestBuilder) GetEmailAppUsageUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*ie49c06630ffed10bc5667592819d2e6ba25594d178187e0c69a98a6cc5ca2e88.GetEmailAppUsageUserDetailWithDateRequestBuilder) {
    return ie49c06630ffed10bc5667592819d2e6ba25594d178187e0c69a98a6cc5ca2e88.NewGetEmailAppUsageUserDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetEmailAppUsageUserDetailWithPeriod provides operations to call the getEmailAppUsageUserDetail method.
func (m *ReportsRequestBuilder) GetEmailAppUsageUserDetailWithPeriod(period *string)(*i285bf21395498d7467b7e1d3f62e43bbbab40d875429e2800566739471aaa586.GetEmailAppUsageUserDetailWithPeriodRequestBuilder) {
    return i285bf21395498d7467b7e1d3f62e43bbbab40d875429e2800566739471aaa586.NewGetEmailAppUsageUserDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetEmailAppUsageVersionsUserCountsWithPeriod provides operations to call the getEmailAppUsageVersionsUserCounts method.
func (m *ReportsRequestBuilder) GetEmailAppUsageVersionsUserCountsWithPeriod(period *string)(*i4ea8142505b6b02804cdd135578592c59e1ebd8c244878c8a27bac9df783863b.GetEmailAppUsageVersionsUserCountsWithPeriodRequestBuilder) {
    return i4ea8142505b6b02804cdd135578592c59e1ebd8c244878c8a27bac9df783863b.NewGetEmailAppUsageVersionsUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetGroupArchivedPrintJobsWithGroupIdWithStartDateTimeWithEndDateTime provides operations to call the getGroupArchivedPrintJobs method.
func (m *ReportsRequestBuilder) GetGroupArchivedPrintJobsWithGroupIdWithStartDateTimeWithEndDateTime(endDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time, groupId *string, startDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)(*i39400581189e192e8c3927b2ac378c799b30d5cb2c7a932a3e2f9802d2a2f95d.GetGroupArchivedPrintJobsWithGroupIdWithStartDateTimeWithEndDateTimeRequestBuilder) {
    return i39400581189e192e8c3927b2ac378c799b30d5cb2c7a932a3e2f9802d2a2f95d.NewGetGroupArchivedPrintJobsWithGroupIdWithStartDateTimeWithEndDateTimeRequestBuilderInternal(m.pathParameters, m.requestAdapter, endDateTime, groupId, startDateTime);
}
// GetM365AppPlatformUserCountsWithPeriod provides operations to call the getM365AppPlatformUserCounts method.
func (m *ReportsRequestBuilder) GetM365AppPlatformUserCountsWithPeriod(period *string)(*ie603b5f7e630bd6d2e83589da6d79d43b190e1d4c7fc28656991ffe48fca4828.GetM365AppPlatformUserCountsWithPeriodRequestBuilder) {
    return ie603b5f7e630bd6d2e83589da6d79d43b190e1d4c7fc28656991ffe48fca4828.NewGetM365AppPlatformUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetM365AppUserCountsWithPeriod provides operations to call the getM365AppUserCounts method.
func (m *ReportsRequestBuilder) GetM365AppUserCountsWithPeriod(period *string)(*ifd36d48788fa091c2ebd9158094174f0f9a871398672c26de206aab3755cc266.GetM365AppUserCountsWithPeriodRequestBuilder) {
    return ifd36d48788fa091c2ebd9158094174f0f9a871398672c26de206aab3755cc266.NewGetM365AppUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetM365AppUserDetailWithDate provides operations to call the getM365AppUserDetail method.
func (m *ReportsRequestBuilder) GetM365AppUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*i753582442e8209033f45d1a367493028141a79145904e5dbd0343dc59d62d8c9.GetM365AppUserDetailWithDateRequestBuilder) {
    return i753582442e8209033f45d1a367493028141a79145904e5dbd0343dc59d62d8c9.NewGetM365AppUserDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetM365AppUserDetailWithPeriod provides operations to call the getM365AppUserDetail method.
func (m *ReportsRequestBuilder) GetM365AppUserDetailWithPeriod(period *string)(*i1d4731cdb4991f39251caf5d5cbcc8c68036c8b25a8c59c7ee8b24f06ba15ead.GetM365AppUserDetailWithPeriodRequestBuilder) {
    return i1d4731cdb4991f39251caf5d5cbcc8c68036c8b25a8c59c7ee8b24f06ba15ead.NewGetM365AppUserDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetMailboxUsageDetailWithPeriod provides operations to call the getMailboxUsageDetail method.
func (m *ReportsRequestBuilder) GetMailboxUsageDetailWithPeriod(period *string)(*i61a78ce1e0b6e2e60e6d50fe39c731a201868c12c59871387bf62c67fa03cc24.GetMailboxUsageDetailWithPeriodRequestBuilder) {
    return i61a78ce1e0b6e2e60e6d50fe39c731a201868c12c59871387bf62c67fa03cc24.NewGetMailboxUsageDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetMailboxUsageMailboxCountsWithPeriod provides operations to call the getMailboxUsageMailboxCounts method.
func (m *ReportsRequestBuilder) GetMailboxUsageMailboxCountsWithPeriod(period *string)(*i700e0ab2f1ea334038b3293b2396e12432ea3545b07ed127bb12d122bbb6bf90.GetMailboxUsageMailboxCountsWithPeriodRequestBuilder) {
    return i700e0ab2f1ea334038b3293b2396e12432ea3545b07ed127bb12d122bbb6bf90.NewGetMailboxUsageMailboxCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetMailboxUsageQuotaStatusMailboxCountsWithPeriod provides operations to call the getMailboxUsageQuotaStatusMailboxCounts method.
func (m *ReportsRequestBuilder) GetMailboxUsageQuotaStatusMailboxCountsWithPeriod(period *string)(*i37d455e922500fc7c5e44df5cc84a103ddbb8470a53860cb0c02e64fd6e8e775.GetMailboxUsageQuotaStatusMailboxCountsWithPeriodRequestBuilder) {
    return i37d455e922500fc7c5e44df5cc84a103ddbb8470a53860cb0c02e64fd6e8e775.NewGetMailboxUsageQuotaStatusMailboxCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetMailboxUsageStorageWithPeriod provides operations to call the getMailboxUsageStorage method.
func (m *ReportsRequestBuilder) GetMailboxUsageStorageWithPeriod(period *string)(*id05b47283d832e55b732685b30d918ee0e303ae639c6447cfff5c7a7089741c7.GetMailboxUsageStorageWithPeriodRequestBuilder) {
    return id05b47283d832e55b732685b30d918ee0e303ae639c6447cfff5c7a7089741c7.NewGetMailboxUsageStorageWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOffice365ActivationCounts provides operations to call the getOffice365ActivationCounts method.
func (m *ReportsRequestBuilder) GetOffice365ActivationCounts()(*i37efe63b21f1e137ba9fdf350d68a656b504d390fe52651da7739ce398ec31c8.GetOffice365ActivationCountsRequestBuilder) {
    return i37efe63b21f1e137ba9fdf350d68a656b504d390fe52651da7739ce398ec31c8.NewGetOffice365ActivationCountsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetOffice365ActivationsUserCounts provides operations to call the getOffice365ActivationsUserCounts method.
func (m *ReportsRequestBuilder) GetOffice365ActivationsUserCounts()(*i5be64ff627afac918209c9e795b4c7ad94cf40584d7376fd0bddf4902f737286.GetOffice365ActivationsUserCountsRequestBuilder) {
    return i5be64ff627afac918209c9e795b4c7ad94cf40584d7376fd0bddf4902f737286.NewGetOffice365ActivationsUserCountsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetOffice365ActivationsUserDetail provides operations to call the getOffice365ActivationsUserDetail method.
func (m *ReportsRequestBuilder) GetOffice365ActivationsUserDetail()(*i67d417cc0a32f095a4d9872596c73a5670bab5ca99f2c39b690a3dc3e53dacb6.GetOffice365ActivationsUserDetailRequestBuilder) {
    return i67d417cc0a32f095a4d9872596c73a5670bab5ca99f2c39b690a3dc3e53dacb6.NewGetOffice365ActivationsUserDetailRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// GetOffice365ActiveUserCountsWithPeriod provides operations to call the getOffice365ActiveUserCounts method.
func (m *ReportsRequestBuilder) GetOffice365ActiveUserCountsWithPeriod(period *string)(*i4ea09032e0c277ad5c699279a1e67097201d217b8d0aac6b5962420831d787bd.GetOffice365ActiveUserCountsWithPeriodRequestBuilder) {
    return i4ea09032e0c277ad5c699279a1e67097201d217b8d0aac6b5962420831d787bd.NewGetOffice365ActiveUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOffice365ActiveUserDetailWithDate provides operations to call the getOffice365ActiveUserDetail method.
func (m *ReportsRequestBuilder) GetOffice365ActiveUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*i0d6612898baf60cf3ad09e5ec9209afdfaffc25ddef2fe3fddf028f8f356d2e1.GetOffice365ActiveUserDetailWithDateRequestBuilder) {
    return i0d6612898baf60cf3ad09e5ec9209afdfaffc25ddef2fe3fddf028f8f356d2e1.NewGetOffice365ActiveUserDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetOffice365ActiveUserDetailWithPeriod provides operations to call the getOffice365ActiveUserDetail method.
func (m *ReportsRequestBuilder) GetOffice365ActiveUserDetailWithPeriod(period *string)(*i6edf15eaf951cb8d179890d37ee652592111a5bfeb9fd5709d3609c7593bed81.GetOffice365ActiveUserDetailWithPeriodRequestBuilder) {
    return i6edf15eaf951cb8d179890d37ee652592111a5bfeb9fd5709d3609c7593bed81.NewGetOffice365ActiveUserDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOffice365GroupsActivityCountsWithPeriod provides operations to call the getOffice365GroupsActivityCounts method.
func (m *ReportsRequestBuilder) GetOffice365GroupsActivityCountsWithPeriod(period *string)(*ifb70744a42eac9e3054301d8328a61f22ab480c9dcaaaa3f398cf2a26e1b242b.GetOffice365GroupsActivityCountsWithPeriodRequestBuilder) {
    return ifb70744a42eac9e3054301d8328a61f22ab480c9dcaaaa3f398cf2a26e1b242b.NewGetOffice365GroupsActivityCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOffice365GroupsActivityDetailWithDate provides operations to call the getOffice365GroupsActivityDetail method.
func (m *ReportsRequestBuilder) GetOffice365GroupsActivityDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*i4319947bf5469fb8e87f0bc08088c383a13726c5be51ca3558867ea6abecc315.GetOffice365GroupsActivityDetailWithDateRequestBuilder) {
    return i4319947bf5469fb8e87f0bc08088c383a13726c5be51ca3558867ea6abecc315.NewGetOffice365GroupsActivityDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetOffice365GroupsActivityDetailWithPeriod provides operations to call the getOffice365GroupsActivityDetail method.
func (m *ReportsRequestBuilder) GetOffice365GroupsActivityDetailWithPeriod(period *string)(*i7774a89bdd2b6bc2f70c1e7c58f236fd1f3aa68a40a72db4898fd9c5d4f66988.GetOffice365GroupsActivityDetailWithPeriodRequestBuilder) {
    return i7774a89bdd2b6bc2f70c1e7c58f236fd1f3aa68a40a72db4898fd9c5d4f66988.NewGetOffice365GroupsActivityDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOffice365GroupsActivityFileCountsWithPeriod provides operations to call the getOffice365GroupsActivityFileCounts method.
func (m *ReportsRequestBuilder) GetOffice365GroupsActivityFileCountsWithPeriod(period *string)(*i21ef6a6d7bc99608d2279d5cc2343f5300ae76d3d739202a0c1bc4ff2d13a391.GetOffice365GroupsActivityFileCountsWithPeriodRequestBuilder) {
    return i21ef6a6d7bc99608d2279d5cc2343f5300ae76d3d739202a0c1bc4ff2d13a391.NewGetOffice365GroupsActivityFileCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOffice365GroupsActivityGroupCountsWithPeriod provides operations to call the getOffice365GroupsActivityGroupCounts method.
func (m *ReportsRequestBuilder) GetOffice365GroupsActivityGroupCountsWithPeriod(period *string)(*i5445cb98613918d3a7a579fddf521b041c55f06714612e2097ed886a37be4286.GetOffice365GroupsActivityGroupCountsWithPeriodRequestBuilder) {
    return i5445cb98613918d3a7a579fddf521b041c55f06714612e2097ed886a37be4286.NewGetOffice365GroupsActivityGroupCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOffice365GroupsActivityStorageWithPeriod provides operations to call the getOffice365GroupsActivityStorage method.
func (m *ReportsRequestBuilder) GetOffice365GroupsActivityStorageWithPeriod(period *string)(*i26d1169a1f986183917395a64e40bb9c1870a54b457336ef64e7efa21c214a15.GetOffice365GroupsActivityStorageWithPeriodRequestBuilder) {
    return i26d1169a1f986183917395a64e40bb9c1870a54b457336ef64e7efa21c214a15.NewGetOffice365GroupsActivityStorageWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOffice365ServicesUserCountsWithPeriod provides operations to call the getOffice365ServicesUserCounts method.
func (m *ReportsRequestBuilder) GetOffice365ServicesUserCountsWithPeriod(period *string)(*i2306eba962c786175114d00783028649abfdfb48f769241b2a7dca15db108361.GetOffice365ServicesUserCountsWithPeriodRequestBuilder) {
    return i2306eba962c786175114d00783028649abfdfb48f769241b2a7dca15db108361.NewGetOffice365ServicesUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOneDriveActivityFileCountsWithPeriod provides operations to call the getOneDriveActivityFileCounts method.
func (m *ReportsRequestBuilder) GetOneDriveActivityFileCountsWithPeriod(period *string)(*iec7256d00a8e37c8dec628d1626f62f1224b45c25e3c24418830b8fede7bc006.GetOneDriveActivityFileCountsWithPeriodRequestBuilder) {
    return iec7256d00a8e37c8dec628d1626f62f1224b45c25e3c24418830b8fede7bc006.NewGetOneDriveActivityFileCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOneDriveActivityUserCountsWithPeriod provides operations to call the getOneDriveActivityUserCounts method.
func (m *ReportsRequestBuilder) GetOneDriveActivityUserCountsWithPeriod(period *string)(*i1ccdf78d8963c250d79813dd0917eaaaee794135fc39e35801b01ea5f2005034.GetOneDriveActivityUserCountsWithPeriodRequestBuilder) {
    return i1ccdf78d8963c250d79813dd0917eaaaee794135fc39e35801b01ea5f2005034.NewGetOneDriveActivityUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOneDriveActivityUserDetailWithDate provides operations to call the getOneDriveActivityUserDetail method.
func (m *ReportsRequestBuilder) GetOneDriveActivityUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*ife6d24553cc77853635c06a75c3587d7c978531be356ad656116b7742f60e1f7.GetOneDriveActivityUserDetailWithDateRequestBuilder) {
    return ife6d24553cc77853635c06a75c3587d7c978531be356ad656116b7742f60e1f7.NewGetOneDriveActivityUserDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetOneDriveActivityUserDetailWithPeriod provides operations to call the getOneDriveActivityUserDetail method.
func (m *ReportsRequestBuilder) GetOneDriveActivityUserDetailWithPeriod(period *string)(*iab84fe97b770990c561344b8be4caaca7e6d962145fd1ad9fc9e851313b38857.GetOneDriveActivityUserDetailWithPeriodRequestBuilder) {
    return iab84fe97b770990c561344b8be4caaca7e6d962145fd1ad9fc9e851313b38857.NewGetOneDriveActivityUserDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOneDriveUsageAccountCountsWithPeriod provides operations to call the getOneDriveUsageAccountCounts method.
func (m *ReportsRequestBuilder) GetOneDriveUsageAccountCountsWithPeriod(period *string)(*i75a65b7f6cf316241684f8aeba4d08fb2c2e04274ef9390b6003ff2a3fe87418.GetOneDriveUsageAccountCountsWithPeriodRequestBuilder) {
    return i75a65b7f6cf316241684f8aeba4d08fb2c2e04274ef9390b6003ff2a3fe87418.NewGetOneDriveUsageAccountCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOneDriveUsageAccountDetailWithDate provides operations to call the getOneDriveUsageAccountDetail method.
func (m *ReportsRequestBuilder) GetOneDriveUsageAccountDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*i01986228cfa34d663610609938a20fe1aa899eb3ae22c44d689bf3a7cc4fcec2.GetOneDriveUsageAccountDetailWithDateRequestBuilder) {
    return i01986228cfa34d663610609938a20fe1aa899eb3ae22c44d689bf3a7cc4fcec2.NewGetOneDriveUsageAccountDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetOneDriveUsageAccountDetailWithPeriod provides operations to call the getOneDriveUsageAccountDetail method.
func (m *ReportsRequestBuilder) GetOneDriveUsageAccountDetailWithPeriod(period *string)(*i316dae3dc0a83b965993ba7f5a9d3b29d2b5dcbd706e8977d4a6cb34548c5275.GetOneDriveUsageAccountDetailWithPeriodRequestBuilder) {
    return i316dae3dc0a83b965993ba7f5a9d3b29d2b5dcbd706e8977d4a6cb34548c5275.NewGetOneDriveUsageAccountDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOneDriveUsageFileCountsWithPeriod provides operations to call the getOneDriveUsageFileCounts method.
func (m *ReportsRequestBuilder) GetOneDriveUsageFileCountsWithPeriod(period *string)(*ibc84667b2a2b840609232b734d0a2873e66315fbe63aeadb3c3d468ad9f84d75.GetOneDriveUsageFileCountsWithPeriodRequestBuilder) {
    return ibc84667b2a2b840609232b734d0a2873e66315fbe63aeadb3c3d468ad9f84d75.NewGetOneDriveUsageFileCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetOneDriveUsageStorageWithPeriod provides operations to call the getOneDriveUsageStorage method.
func (m *ReportsRequestBuilder) GetOneDriveUsageStorageWithPeriod(period *string)(*if97af71691601e0f0046051c8e82fb71f8db2f1d3b5f8e7e40346d3fa47b196f.GetOneDriveUsageStorageWithPeriodRequestBuilder) {
    return if97af71691601e0f0046051c8e82fb71f8db2f1d3b5f8e7e40346d3fa47b196f.NewGetOneDriveUsageStorageWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetPrinterArchivedPrintJobsWithPrinterIdWithStartDateTimeWithEndDateTime provides operations to call the getPrinterArchivedPrintJobs method.
func (m *ReportsRequestBuilder) GetPrinterArchivedPrintJobsWithPrinterIdWithStartDateTimeWithEndDateTime(endDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time, printerId *string, startDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)(*i3ed9439789e0857666ebca3cd16359b4745ffe979a682ea73e6cd26f8054ae21.GetPrinterArchivedPrintJobsWithPrinterIdWithStartDateTimeWithEndDateTimeRequestBuilder) {
    return i3ed9439789e0857666ebca3cd16359b4745ffe979a682ea73e6cd26f8054ae21.NewGetPrinterArchivedPrintJobsWithPrinterIdWithStartDateTimeWithEndDateTimeRequestBuilderInternal(m.pathParameters, m.requestAdapter, endDateTime, printerId, startDateTime);
}
// GetSharePointActivityFileCountsWithPeriod provides operations to call the getSharePointActivityFileCounts method.
func (m *ReportsRequestBuilder) GetSharePointActivityFileCountsWithPeriod(period *string)(*i1bfeeda14fdcd3c4bc1c9b02a4d925837c77654498184deaffbfeb8f496c10c5.GetSharePointActivityFileCountsWithPeriodRequestBuilder) {
    return i1bfeeda14fdcd3c4bc1c9b02a4d925837c77654498184deaffbfeb8f496c10c5.NewGetSharePointActivityFileCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSharePointActivityPagesWithPeriod provides operations to call the getSharePointActivityPages method.
func (m *ReportsRequestBuilder) GetSharePointActivityPagesWithPeriod(period *string)(*i6b3e5d0662b4c04bf1f969a3c62a48d50c4f7ecb57a3913782ce351d87a0035e.GetSharePointActivityPagesWithPeriodRequestBuilder) {
    return i6b3e5d0662b4c04bf1f969a3c62a48d50c4f7ecb57a3913782ce351d87a0035e.NewGetSharePointActivityPagesWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSharePointActivityUserCountsWithPeriod provides operations to call the getSharePointActivityUserCounts method.
func (m *ReportsRequestBuilder) GetSharePointActivityUserCountsWithPeriod(period *string)(*i40f65a5b54f01dfaea303a3582e82bcdaa71c2195cc00b7915431e77f38e9a0f.GetSharePointActivityUserCountsWithPeriodRequestBuilder) {
    return i40f65a5b54f01dfaea303a3582e82bcdaa71c2195cc00b7915431e77f38e9a0f.NewGetSharePointActivityUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSharePointActivityUserDetailWithDate provides operations to call the getSharePointActivityUserDetail method.
func (m *ReportsRequestBuilder) GetSharePointActivityUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*ib3bef6d1d1d843d8e20e58b41c29e9b82bb511ac97e381039f942ffee79a3b9a.GetSharePointActivityUserDetailWithDateRequestBuilder) {
    return ib3bef6d1d1d843d8e20e58b41c29e9b82bb511ac97e381039f942ffee79a3b9a.NewGetSharePointActivityUserDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetSharePointActivityUserDetailWithPeriod provides operations to call the getSharePointActivityUserDetail method.
func (m *ReportsRequestBuilder) GetSharePointActivityUserDetailWithPeriod(period *string)(*i56c28dc52ab4ecbfe37baebbfaf6f279f8d130872c0931a194264028d89e2c96.GetSharePointActivityUserDetailWithPeriodRequestBuilder) {
    return i56c28dc52ab4ecbfe37baebbfaf6f279f8d130872c0931a194264028d89e2c96.NewGetSharePointActivityUserDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSharePointSiteUsageDetailWithDate provides operations to call the getSharePointSiteUsageDetail method.
func (m *ReportsRequestBuilder) GetSharePointSiteUsageDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*ic5ac4372354eec289834ac2ae20a4b14244024792aa3c3ad01b0f0a204b1e81b.GetSharePointSiteUsageDetailWithDateRequestBuilder) {
    return ic5ac4372354eec289834ac2ae20a4b14244024792aa3c3ad01b0f0a204b1e81b.NewGetSharePointSiteUsageDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetSharePointSiteUsageDetailWithPeriod provides operations to call the getSharePointSiteUsageDetail method.
func (m *ReportsRequestBuilder) GetSharePointSiteUsageDetailWithPeriod(period *string)(*ib9c6325caeea63e2dd57f99af77223cbb370ce799b27f8b8a3f6a3f987a9f9fe.GetSharePointSiteUsageDetailWithPeriodRequestBuilder) {
    return ib9c6325caeea63e2dd57f99af77223cbb370ce799b27f8b8a3f6a3f987a9f9fe.NewGetSharePointSiteUsageDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSharePointSiteUsageFileCountsWithPeriod provides operations to call the getSharePointSiteUsageFileCounts method.
func (m *ReportsRequestBuilder) GetSharePointSiteUsageFileCountsWithPeriod(period *string)(*i2699e05c2942728d3c9dbbec2f7da6269f18ec269bdc96532cb49ca45ea2e03f.GetSharePointSiteUsageFileCountsWithPeriodRequestBuilder) {
    return i2699e05c2942728d3c9dbbec2f7da6269f18ec269bdc96532cb49ca45ea2e03f.NewGetSharePointSiteUsageFileCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSharePointSiteUsagePagesWithPeriod provides operations to call the getSharePointSiteUsagePages method.
func (m *ReportsRequestBuilder) GetSharePointSiteUsagePagesWithPeriod(period *string)(*i2f04ee022b2f4dd0c31e5c7f89cc28c7d1e36585beb573c9c689634579fb8999.GetSharePointSiteUsagePagesWithPeriodRequestBuilder) {
    return i2f04ee022b2f4dd0c31e5c7f89cc28c7d1e36585beb573c9c689634579fb8999.NewGetSharePointSiteUsagePagesWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSharePointSiteUsageSiteCountsWithPeriod provides operations to call the getSharePointSiteUsageSiteCounts method.
func (m *ReportsRequestBuilder) GetSharePointSiteUsageSiteCountsWithPeriod(period *string)(*i347d730506a06c407296d7557c653293a35a9c38b50c7256c1c7abb4c736ef9a.GetSharePointSiteUsageSiteCountsWithPeriodRequestBuilder) {
    return i347d730506a06c407296d7557c653293a35a9c38b50c7256c1c7abb4c736ef9a.NewGetSharePointSiteUsageSiteCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSharePointSiteUsageStorageWithPeriod provides operations to call the getSharePointSiteUsageStorage method.
func (m *ReportsRequestBuilder) GetSharePointSiteUsageStorageWithPeriod(period *string)(*i028952a5521bf457c8a8796489c7f5a1137cf7b0745234478ff5fd12bb087564.GetSharePointSiteUsageStorageWithPeriodRequestBuilder) {
    return i028952a5521bf457c8a8796489c7f5a1137cf7b0745234478ff5fd12bb087564.NewGetSharePointSiteUsageStorageWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessActivityCountsWithPeriod provides operations to call the getSkypeForBusinessActivityCounts method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessActivityCountsWithPeriod(period *string)(*i8ad77c0341e431729102d8b93f9066b82a7f4fc8ee3e748eeea007bf3d9d5b18.GetSkypeForBusinessActivityCountsWithPeriodRequestBuilder) {
    return i8ad77c0341e431729102d8b93f9066b82a7f4fc8ee3e748eeea007bf3d9d5b18.NewGetSkypeForBusinessActivityCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessActivityUserCountsWithPeriod provides operations to call the getSkypeForBusinessActivityUserCounts method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessActivityUserCountsWithPeriod(period *string)(*i069326085999a3926845312c56f68613c15a7dcb468bf7f0592670ea8f5896aa.GetSkypeForBusinessActivityUserCountsWithPeriodRequestBuilder) {
    return i069326085999a3926845312c56f68613c15a7dcb468bf7f0592670ea8f5896aa.NewGetSkypeForBusinessActivityUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessActivityUserDetailWithDate provides operations to call the getSkypeForBusinessActivityUserDetail method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessActivityUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*i6a26bcea122cbd3594cacac16b0715042c3ca45bb5f889d1fd0bd83d243a6fe3.GetSkypeForBusinessActivityUserDetailWithDateRequestBuilder) {
    return i6a26bcea122cbd3594cacac16b0715042c3ca45bb5f889d1fd0bd83d243a6fe3.NewGetSkypeForBusinessActivityUserDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetSkypeForBusinessActivityUserDetailWithPeriod provides operations to call the getSkypeForBusinessActivityUserDetail method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessActivityUserDetailWithPeriod(period *string)(*ia79ac9c43dc0218a9c3f469a247120979433dd0e3e53de3949570224b3cf3936.GetSkypeForBusinessActivityUserDetailWithPeriodRequestBuilder) {
    return ia79ac9c43dc0218a9c3f469a247120979433dd0e3e53de3949570224b3cf3936.NewGetSkypeForBusinessActivityUserDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessDeviceUsageDistributionUserCountsWithPeriod provides operations to call the getSkypeForBusinessDeviceUsageDistributionUserCounts method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessDeviceUsageDistributionUserCountsWithPeriod(period *string)(*ia5f41691c1fbeb62b6b1de131fa27f7306db515c8759a9e49ec88553c0e6d97d.GetSkypeForBusinessDeviceUsageDistributionUserCountsWithPeriodRequestBuilder) {
    return ia5f41691c1fbeb62b6b1de131fa27f7306db515c8759a9e49ec88553c0e6d97d.NewGetSkypeForBusinessDeviceUsageDistributionUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessDeviceUsageUserCountsWithPeriod provides operations to call the getSkypeForBusinessDeviceUsageUserCounts method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessDeviceUsageUserCountsWithPeriod(period *string)(*i9de2dd5236ac703b10206d8522e5a8f5490288eadb62a3fbf88d9becd9f061bf.GetSkypeForBusinessDeviceUsageUserCountsWithPeriodRequestBuilder) {
    return i9de2dd5236ac703b10206d8522e5a8f5490288eadb62a3fbf88d9becd9f061bf.NewGetSkypeForBusinessDeviceUsageUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessDeviceUsageUserDetailWithDate provides operations to call the getSkypeForBusinessDeviceUsageUserDetail method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessDeviceUsageUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*i037933f7467ae74beac89af9fadc6cf1a5a8272d8f21a5c2f98bc8f8e7c28ab6.GetSkypeForBusinessDeviceUsageUserDetailWithDateRequestBuilder) {
    return i037933f7467ae74beac89af9fadc6cf1a5a8272d8f21a5c2f98bc8f8e7c28ab6.NewGetSkypeForBusinessDeviceUsageUserDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetSkypeForBusinessDeviceUsageUserDetailWithPeriod provides operations to call the getSkypeForBusinessDeviceUsageUserDetail method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessDeviceUsageUserDetailWithPeriod(period *string)(*i6c58964b36e86bf47df5ddf5c8161ca0e4621aee61aaa9517f1dded3120cfc8d.GetSkypeForBusinessDeviceUsageUserDetailWithPeriodRequestBuilder) {
    return i6c58964b36e86bf47df5ddf5c8161ca0e4621aee61aaa9517f1dded3120cfc8d.NewGetSkypeForBusinessDeviceUsageUserDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessOrganizerActivityCountsWithPeriod provides operations to call the getSkypeForBusinessOrganizerActivityCounts method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessOrganizerActivityCountsWithPeriod(period *string)(*i2e849a1f15987664ba1a6648f0c1d0eb9e2ca5965a632cc3761ad14bb91674bf.GetSkypeForBusinessOrganizerActivityCountsWithPeriodRequestBuilder) {
    return i2e849a1f15987664ba1a6648f0c1d0eb9e2ca5965a632cc3761ad14bb91674bf.NewGetSkypeForBusinessOrganizerActivityCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessOrganizerActivityMinuteCountsWithPeriod provides operations to call the getSkypeForBusinessOrganizerActivityMinuteCounts method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessOrganizerActivityMinuteCountsWithPeriod(period *string)(*i8d02e7a0b578c5b4465f96608d724ae82f61e67266f11ff2ab8aad96d55d1ac1.GetSkypeForBusinessOrganizerActivityMinuteCountsWithPeriodRequestBuilder) {
    return i8d02e7a0b578c5b4465f96608d724ae82f61e67266f11ff2ab8aad96d55d1ac1.NewGetSkypeForBusinessOrganizerActivityMinuteCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessOrganizerActivityUserCountsWithPeriod provides operations to call the getSkypeForBusinessOrganizerActivityUserCounts method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessOrganizerActivityUserCountsWithPeriod(period *string)(*i066e5ea9f7e2ed1b047023128aab6caef0db934d50a71df44b7b6cf7bfb5a10e.GetSkypeForBusinessOrganizerActivityUserCountsWithPeriodRequestBuilder) {
    return i066e5ea9f7e2ed1b047023128aab6caef0db934d50a71df44b7b6cf7bfb5a10e.NewGetSkypeForBusinessOrganizerActivityUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessParticipantActivityCountsWithPeriod provides operations to call the getSkypeForBusinessParticipantActivityCounts method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessParticipantActivityCountsWithPeriod(period *string)(*i6e8f025e6b5802865cf6c583d64d93f7e1c603ead045ca53bef1f6b9cde8a5b0.GetSkypeForBusinessParticipantActivityCountsWithPeriodRequestBuilder) {
    return i6e8f025e6b5802865cf6c583d64d93f7e1c603ead045ca53bef1f6b9cde8a5b0.NewGetSkypeForBusinessParticipantActivityCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessParticipantActivityMinuteCountsWithPeriod provides operations to call the getSkypeForBusinessParticipantActivityMinuteCounts method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessParticipantActivityMinuteCountsWithPeriod(period *string)(*idcca069eb97aac5c90ca4ddb83fee2301a14619967cfdadbb5d0a81effc5cc59.GetSkypeForBusinessParticipantActivityMinuteCountsWithPeriodRequestBuilder) {
    return idcca069eb97aac5c90ca4ddb83fee2301a14619967cfdadbb5d0a81effc5cc59.NewGetSkypeForBusinessParticipantActivityMinuteCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessParticipantActivityUserCountsWithPeriod provides operations to call the getSkypeForBusinessParticipantActivityUserCounts method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessParticipantActivityUserCountsWithPeriod(period *string)(*i05ca24f1e64f29c556c3c67a9024eea8b3aac4d350f05ff2c52f4511bf0b6f33.GetSkypeForBusinessParticipantActivityUserCountsWithPeriodRequestBuilder) {
    return i05ca24f1e64f29c556c3c67a9024eea8b3aac4d350f05ff2c52f4511bf0b6f33.NewGetSkypeForBusinessParticipantActivityUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessPeerToPeerActivityCountsWithPeriod provides operations to call the getSkypeForBusinessPeerToPeerActivityCounts method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessPeerToPeerActivityCountsWithPeriod(period *string)(*i200c880280f3a34037634650d2cd9969f9a07037f773af21ebb6141e40dde23b.GetSkypeForBusinessPeerToPeerActivityCountsWithPeriodRequestBuilder) {
    return i200c880280f3a34037634650d2cd9969f9a07037f773af21ebb6141e40dde23b.NewGetSkypeForBusinessPeerToPeerActivityCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessPeerToPeerActivityMinuteCountsWithPeriod provides operations to call the getSkypeForBusinessPeerToPeerActivityMinuteCounts method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessPeerToPeerActivityMinuteCountsWithPeriod(period *string)(*ic5508dd2223dfecdf5814372064e28b262f8f8bfe997c49795b7f76e30889837.GetSkypeForBusinessPeerToPeerActivityMinuteCountsWithPeriodRequestBuilder) {
    return ic5508dd2223dfecdf5814372064e28b262f8f8bfe997c49795b7f76e30889837.NewGetSkypeForBusinessPeerToPeerActivityMinuteCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetSkypeForBusinessPeerToPeerActivityUserCountsWithPeriod provides operations to call the getSkypeForBusinessPeerToPeerActivityUserCounts method.
func (m *ReportsRequestBuilder) GetSkypeForBusinessPeerToPeerActivityUserCountsWithPeriod(period *string)(*i416e8a1993d00c5e5b47ab655f48ac49786eb1442a5459b7cb7023ac34c28093.GetSkypeForBusinessPeerToPeerActivityUserCountsWithPeriodRequestBuilder) {
    return i416e8a1993d00c5e5b47ab655f48ac49786eb1442a5459b7cb7023ac34c28093.NewGetSkypeForBusinessPeerToPeerActivityUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetTeamsDeviceUsageDistributionUserCountsWithPeriod provides operations to call the getTeamsDeviceUsageDistributionUserCounts method.
func (m *ReportsRequestBuilder) GetTeamsDeviceUsageDistributionUserCountsWithPeriod(period *string)(*i05c8d276b86beee65c00c115937a46c2753c7d328e142fe9afab6ef90b5a1a79.GetTeamsDeviceUsageDistributionUserCountsWithPeriodRequestBuilder) {
    return i05c8d276b86beee65c00c115937a46c2753c7d328e142fe9afab6ef90b5a1a79.NewGetTeamsDeviceUsageDistributionUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetTeamsDeviceUsageUserCountsWithPeriod provides operations to call the getTeamsDeviceUsageUserCounts method.
func (m *ReportsRequestBuilder) GetTeamsDeviceUsageUserCountsWithPeriod(period *string)(*ife99981e20f3609c2111e3e2aff666b997075ba1e7c98ff44f1950f1d6265ccb.GetTeamsDeviceUsageUserCountsWithPeriodRequestBuilder) {
    return ife99981e20f3609c2111e3e2aff666b997075ba1e7c98ff44f1950f1d6265ccb.NewGetTeamsDeviceUsageUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetTeamsDeviceUsageUserDetailWithDate provides operations to call the getTeamsDeviceUsageUserDetail method.
func (m *ReportsRequestBuilder) GetTeamsDeviceUsageUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*i677bcf76c9544b6088f18b0d84351cd554f0537738c43ca67e647fffec66978c.GetTeamsDeviceUsageUserDetailWithDateRequestBuilder) {
    return i677bcf76c9544b6088f18b0d84351cd554f0537738c43ca67e647fffec66978c.NewGetTeamsDeviceUsageUserDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetTeamsDeviceUsageUserDetailWithPeriod provides operations to call the getTeamsDeviceUsageUserDetail method.
func (m *ReportsRequestBuilder) GetTeamsDeviceUsageUserDetailWithPeriod(period *string)(*iaaacf2b7aacf4d5f8107fa18d6f0d6cb33066e3335d8ab8cc1a2adfa1f4e3c42.GetTeamsDeviceUsageUserDetailWithPeriodRequestBuilder) {
    return iaaacf2b7aacf4d5f8107fa18d6f0d6cb33066e3335d8ab8cc1a2adfa1f4e3c42.NewGetTeamsDeviceUsageUserDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetTeamsUserActivityCountsWithPeriod provides operations to call the getTeamsUserActivityCounts method.
func (m *ReportsRequestBuilder) GetTeamsUserActivityCountsWithPeriod(period *string)(*i9b65430927e689a26a700312296e423fc6e6e6c0bb86f72287b75c7915e695e0.GetTeamsUserActivityCountsWithPeriodRequestBuilder) {
    return i9b65430927e689a26a700312296e423fc6e6e6c0bb86f72287b75c7915e695e0.NewGetTeamsUserActivityCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetTeamsUserActivityUserCountsWithPeriod provides operations to call the getTeamsUserActivityUserCounts method.
func (m *ReportsRequestBuilder) GetTeamsUserActivityUserCountsWithPeriod(period *string)(*i34861f91c890c43257a1e1386ba24de1ffe44f842eadf95d13de30d27a852991.GetTeamsUserActivityUserCountsWithPeriodRequestBuilder) {
    return i34861f91c890c43257a1e1386ba24de1ffe44f842eadf95d13de30d27a852991.NewGetTeamsUserActivityUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetTeamsUserActivityUserDetailWithDate provides operations to call the getTeamsUserActivityUserDetail method.
func (m *ReportsRequestBuilder) GetTeamsUserActivityUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*i0c4b2ea99f0789f3909bf8d293ade8bd9216c7010a8e0dc72ceea54c21c1d481.GetTeamsUserActivityUserDetailWithDateRequestBuilder) {
    return i0c4b2ea99f0789f3909bf8d293ade8bd9216c7010a8e0dc72ceea54c21c1d481.NewGetTeamsUserActivityUserDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetTeamsUserActivityUserDetailWithPeriod provides operations to call the getTeamsUserActivityUserDetail method.
func (m *ReportsRequestBuilder) GetTeamsUserActivityUserDetailWithPeriod(period *string)(*ibfd05176830665cf1dae7366e8b491572f134d04e12120a5499c283d6c8ae062.GetTeamsUserActivityUserDetailWithPeriodRequestBuilder) {
    return ibfd05176830665cf1dae7366e8b491572f134d04e12120a5499c283d6c8ae062.NewGetTeamsUserActivityUserDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetUserArchivedPrintJobsWithUserIdWithStartDateTimeWithEndDateTime provides operations to call the getUserArchivedPrintJobs method.
func (m *ReportsRequestBuilder) GetUserArchivedPrintJobsWithUserIdWithStartDateTimeWithEndDateTime(endDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time, startDateTime *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time, userId *string)(*iafadefc9c12601a3611573bb976b73b251f9ea476a68e3bdc9ceab9e0d5cc308.GetUserArchivedPrintJobsWithUserIdWithStartDateTimeWithEndDateTimeRequestBuilder) {
    return iafadefc9c12601a3611573bb976b73b251f9ea476a68e3bdc9ceab9e0d5cc308.NewGetUserArchivedPrintJobsWithUserIdWithStartDateTimeWithEndDateTimeRequestBuilderInternal(m.pathParameters, m.requestAdapter, endDateTime, startDateTime, userId);
}
// GetYammerActivityCountsWithPeriod provides operations to call the getYammerActivityCounts method.
func (m *ReportsRequestBuilder) GetYammerActivityCountsWithPeriod(period *string)(*ifd45221e0ff465bb6f37fa740578c697634181cf333e3c01e4db49019e2ba5eb.GetYammerActivityCountsWithPeriodRequestBuilder) {
    return ifd45221e0ff465bb6f37fa740578c697634181cf333e3c01e4db49019e2ba5eb.NewGetYammerActivityCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetYammerActivityUserCountsWithPeriod provides operations to call the getYammerActivityUserCounts method.
func (m *ReportsRequestBuilder) GetYammerActivityUserCountsWithPeriod(period *string)(*i13fef70f90226be585263be7022fe29fa384901ba2188c080017c9272ca23e99.GetYammerActivityUserCountsWithPeriodRequestBuilder) {
    return i13fef70f90226be585263be7022fe29fa384901ba2188c080017c9272ca23e99.NewGetYammerActivityUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetYammerActivityUserDetailWithDate provides operations to call the getYammerActivityUserDetail method.
func (m *ReportsRequestBuilder) GetYammerActivityUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*id7d2a94ac1ea97cc24d2cd9cb47182b168535107d0d09ad275c4e955f417e029.GetYammerActivityUserDetailWithDateRequestBuilder) {
    return id7d2a94ac1ea97cc24d2cd9cb47182b168535107d0d09ad275c4e955f417e029.NewGetYammerActivityUserDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetYammerActivityUserDetailWithPeriod provides operations to call the getYammerActivityUserDetail method.
func (m *ReportsRequestBuilder) GetYammerActivityUserDetailWithPeriod(period *string)(*id379a8afbc53e3411a43fbddcd71ea21f90a2a71ada609569a5c9040ebee5f51.GetYammerActivityUserDetailWithPeriodRequestBuilder) {
    return id379a8afbc53e3411a43fbddcd71ea21f90a2a71ada609569a5c9040ebee5f51.NewGetYammerActivityUserDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetYammerDeviceUsageDistributionUserCountsWithPeriod provides operations to call the getYammerDeviceUsageDistributionUserCounts method.
func (m *ReportsRequestBuilder) GetYammerDeviceUsageDistributionUserCountsWithPeriod(period *string)(*i03d33a99775bdfdb348179d972066b208c7410da4d93f4776c48bb0c7619f539.GetYammerDeviceUsageDistributionUserCountsWithPeriodRequestBuilder) {
    return i03d33a99775bdfdb348179d972066b208c7410da4d93f4776c48bb0c7619f539.NewGetYammerDeviceUsageDistributionUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetYammerDeviceUsageUserCountsWithPeriod provides operations to call the getYammerDeviceUsageUserCounts method.
func (m *ReportsRequestBuilder) GetYammerDeviceUsageUserCountsWithPeriod(period *string)(*i9b77e3c3d15971fe31fa9e23d6bdc4c74a584c922585117306c4db805ac4131c.GetYammerDeviceUsageUserCountsWithPeriodRequestBuilder) {
    return i9b77e3c3d15971fe31fa9e23d6bdc4c74a584c922585117306c4db805ac4131c.NewGetYammerDeviceUsageUserCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetYammerDeviceUsageUserDetailWithDate provides operations to call the getYammerDeviceUsageUserDetail method.
func (m *ReportsRequestBuilder) GetYammerDeviceUsageUserDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*if97f8cb01f2d6a02697169f75d571d9a67fc56f2041590c23f83bf2114377bda.GetYammerDeviceUsageUserDetailWithDateRequestBuilder) {
    return if97f8cb01f2d6a02697169f75d571d9a67fc56f2041590c23f83bf2114377bda.NewGetYammerDeviceUsageUserDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetYammerDeviceUsageUserDetailWithPeriod provides operations to call the getYammerDeviceUsageUserDetail method.
func (m *ReportsRequestBuilder) GetYammerDeviceUsageUserDetailWithPeriod(period *string)(*i6e07bd3dadcc8aa2d52f07e77f1f35672ac7ce62aa2917c7a9e4e77b8f4c4db5.GetYammerDeviceUsageUserDetailWithPeriodRequestBuilder) {
    return i6e07bd3dadcc8aa2d52f07e77f1f35672ac7ce62aa2917c7a9e4e77b8f4c4db5.NewGetYammerDeviceUsageUserDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetYammerGroupsActivityCountsWithPeriod provides operations to call the getYammerGroupsActivityCounts method.
func (m *ReportsRequestBuilder) GetYammerGroupsActivityCountsWithPeriod(period *string)(*i7ee204dcbd579a0c0d252ddbfcdd01bf9f244b2b38f671f02c4108f867a99c73.GetYammerGroupsActivityCountsWithPeriodRequestBuilder) {
    return i7ee204dcbd579a0c0d252ddbfcdd01bf9f244b2b38f671f02c4108f867a99c73.NewGetYammerGroupsActivityCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetYammerGroupsActivityDetailWithDate provides operations to call the getYammerGroupsActivityDetail method.
func (m *ReportsRequestBuilder) GetYammerGroupsActivityDetailWithDate(date *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.DateOnly)(*ia3da511007b4da9fa7b31d42bcef05972fa35df25142f53ebb0f2944e3932886.GetYammerGroupsActivityDetailWithDateRequestBuilder) {
    return ia3da511007b4da9fa7b31d42bcef05972fa35df25142f53ebb0f2944e3932886.NewGetYammerGroupsActivityDetailWithDateRequestBuilderInternal(m.pathParameters, m.requestAdapter, date);
}
// GetYammerGroupsActivityDetailWithPeriod provides operations to call the getYammerGroupsActivityDetail method.
func (m *ReportsRequestBuilder) GetYammerGroupsActivityDetailWithPeriod(period *string)(*i5277312ad2cbf60385fe9448b03ae55d7aefafa14a9c09b135b745563e0dd271.GetYammerGroupsActivityDetailWithPeriodRequestBuilder) {
    return i5277312ad2cbf60385fe9448b03ae55d7aefafa14a9c09b135b745563e0dd271.NewGetYammerGroupsActivityDetailWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// GetYammerGroupsActivityGroupCountsWithPeriod provides operations to call the getYammerGroupsActivityGroupCounts method.
func (m *ReportsRequestBuilder) GetYammerGroupsActivityGroupCountsWithPeriod(period *string)(*i55b749bfea07392ccbd597dc8d76a153720e39963b1e66527cac5bb79f8c6b29.GetYammerGroupsActivityGroupCountsWithPeriodRequestBuilder) {
    return i55b749bfea07392ccbd597dc8d76a153720e39963b1e66527cac5bb79f8c6b29.NewGetYammerGroupsActivityGroupCountsWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// ManagedDeviceEnrollmentFailureDetails provides operations to call the managedDeviceEnrollmentFailureDetails method.
func (m *ReportsRequestBuilder) ManagedDeviceEnrollmentFailureDetails()(*i0e6a5a79552b1c4826c9418d405e7420746bbd08e606bfd5c84ba1ad674aed8f.ManagedDeviceEnrollmentFailureDetailsRequestBuilder) {
    return i0e6a5a79552b1c4826c9418d405e7420746bbd08e606bfd5c84ba1ad674aed8f.NewManagedDeviceEnrollmentFailureDetailsRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ManagedDeviceEnrollmentFailureDetailsWithSkipWithTopWithFilterWithSkipToken provides operations to call the managedDeviceEnrollmentFailureDetails method.
func (m *ReportsRequestBuilder) ManagedDeviceEnrollmentFailureDetailsWithSkipWithTopWithFilterWithSkipToken(filter *string, skip *int32, skipToken *string, top *int32)(*i374eb5f5ea7befd0f519a8a1c2886282c8ae9910332bd6b507ad8b29eddfdf2d.ManagedDeviceEnrollmentFailureDetailsWithSkipWithTopWithFilterWithSkipTokenRequestBuilder) {
    return i374eb5f5ea7befd0f519a8a1c2886282c8ae9910332bd6b507ad8b29eddfdf2d.NewManagedDeviceEnrollmentFailureDetailsWithSkipWithTopWithFilterWithSkipTokenRequestBuilderInternal(m.pathParameters, m.requestAdapter, filter, skip, skipToken, top);
}
// ManagedDeviceEnrollmentTopFailures provides operations to call the managedDeviceEnrollmentTopFailures method.
func (m *ReportsRequestBuilder) ManagedDeviceEnrollmentTopFailures()(*i9d1992f642df096181dd0b243f5639c29682f838689705077b65251f267d222f.ManagedDeviceEnrollmentTopFailuresRequestBuilder) {
    return i9d1992f642df096181dd0b243f5639c29682f838689705077b65251f267d222f.NewManagedDeviceEnrollmentTopFailuresRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// ManagedDeviceEnrollmentTopFailuresWithPeriod provides operations to call the managedDeviceEnrollmentTopFailures method.
func (m *ReportsRequestBuilder) ManagedDeviceEnrollmentTopFailuresWithPeriod(period *string)(*i3f4aebd4163df4cb0a4f38097f15a9f8ba73853688cf5a50b93a222eb4832abe.ManagedDeviceEnrollmentTopFailuresWithPeriodRequestBuilder) {
    return i3f4aebd4163df4cb0a4f38097f15a9f8ba73853688cf5a50b93a222eb4832abe.NewManagedDeviceEnrollmentTopFailuresWithPeriodRequestBuilderInternal(m.pathParameters, m.requestAdapter, period);
}
// MonthlyPrintUsageByPrinter provides operations to manage the monthlyPrintUsageByPrinter property of the microsoft.graph.reportRoot entity.
func (m *ReportsRequestBuilder) MonthlyPrintUsageByPrinter()(*if1205fb664936127831759fefe4c2c43852474a3fc0419f51e89975f0f281d32.MonthlyPrintUsageByPrinterRequestBuilder) {
    return if1205fb664936127831759fefe4c2c43852474a3fc0419f51e89975f0f281d32.NewMonthlyPrintUsageByPrinterRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MonthlyPrintUsageByPrinterById provides operations to manage the monthlyPrintUsageByPrinter property of the microsoft.graph.reportRoot entity.
func (m *ReportsRequestBuilder) MonthlyPrintUsageByPrinterById(id string)(*i0fbd1879eec3f09756dabd53c9de402f839fcf5ed7e8be0da0d3e4135b124b0a.PrintUsageByPrinterItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["printUsageByPrinter%2Did"] = id
    }
    return i0fbd1879eec3f09756dabd53c9de402f839fcf5ed7e8be0da0d3e4135b124b0a.NewPrintUsageByPrinterItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// MonthlyPrintUsageByUser provides operations to manage the monthlyPrintUsageByUser property of the microsoft.graph.reportRoot entity.
func (m *ReportsRequestBuilder) MonthlyPrintUsageByUser()(*i3f12902a7f7457a60f97014c6db3b2bd89dac8cabe4d467b766677db3d6029d0.MonthlyPrintUsageByUserRequestBuilder) {
    return i3f12902a7f7457a60f97014c6db3b2bd89dac8cabe4d467b766677db3d6029d0.NewMonthlyPrintUsageByUserRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
// MonthlyPrintUsageByUserById provides operations to manage the monthlyPrintUsageByUser property of the microsoft.graph.reportRoot entity.
func (m *ReportsRequestBuilder) MonthlyPrintUsageByUserById(id string)(*if4346a9921b148666613b317b233f49301e4eef2f130f6b9e02414250e5a2294.PrintUsageByUserItemRequestBuilder) {
    urlTplParams := make(map[string]string)
    for idx, item := range m.pathParameters {
        urlTplParams[idx] = item
    }
    if id != "" {
        urlTplParams["printUsageByUser%2Did"] = id
    }
    return if4346a9921b148666613b317b233f49301e4eef2f130f6b9e02414250e5a2294.NewPrintUsageByUserItemRequestBuilderInternal(urlTplParams, m.requestAdapter);
}
// Patch update reports
func (m *ReportsRequestBuilder) Patch(ctx context.Context, body iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ReportRootable, requestConfiguration *ReportsRequestBuilderPatchRequestConfiguration)(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ReportRootable, error) {
    requestInfo, err := m.CreatePatchRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
        "5XX": ia572726a95efa92ddd544552cd950653dc691023836923576b2f4bf716cf204a.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateReportRootFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.ReportRootable), nil
}
// Security provides operations to manage the security property of the microsoft.graph.reportRoot entity.
func (m *ReportsRequestBuilder) Security()(*i2ee11973211f7fc9312ca00d777ef8304aada9d21bbf7c226dda5bef84678886.SecurityRequestBuilder) {
    return i2ee11973211f7fc9312ca00d777ef8304aada9d21bbf7c226dda5bef84678886.NewSecurityRequestBuilderInternal(m.pathParameters, m.requestAdapter);
}
