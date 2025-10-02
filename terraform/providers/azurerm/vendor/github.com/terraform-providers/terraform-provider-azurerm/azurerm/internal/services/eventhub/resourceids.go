package eventhub

//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Cluster -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/clusters/cluster1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=EventHubConsumerGroup -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/eventhubs/eventhub1/consumergroups/consumergroup1
//go:generate go run ../../tools/generator-resource-id/main.go -path=./ -name=Namespace -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1
//go:generate go run ../../tools/generator-resource-id/main.go -rewrite=true -path=./ -name=NamespaceAuthorizationRule -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1/authorizationRules/rule1
