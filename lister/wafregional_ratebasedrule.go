package lister

import (
	"fmt"

	"github.com/trek10inc/awsets/context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafregional"
	"github.com/trek10inc/awsets/resource"
)

type AWSWafRegionalRateBasedRule struct {
}

func init() {
	i := AWSWafRegionalRateBasedRule{}
	listers = append(listers, i)
}

func (l AWSWafRegionalRateBasedRule) Types() []resource.ResourceType {
	return []resource.ResourceType{resource.WafRegionalRateBasedRule}
}

func (l AWSWafRegionalRateBasedRule) List(ctx context.AWSetsCtx) (*resource.Group, error) {
	svc := wafregional.New(ctx.AWSCfg)
	rg := resource.NewGroup()

	var nextMarker *string
	for {
		res, err := svc.ListRateBasedRulesRequest(&wafregional.ListRateBasedRulesInput{
			Limit:      aws.Int64(100),
			NextMarker: nextMarker,
		}).Send(ctx.Context)
		if err != nil {
			return rg, fmt.Errorf("failed to list regional rate based rule: %w", err)
		}
		if res.ListRateBasedRulesOutput == nil {
			continue
		}
		for _, id := range res.ListRateBasedRulesOutput.Rules {
			rule, err := svc.GetRateBasedRuleRequest(&wafregional.GetRateBasedRuleInput{
				RuleId: id.RuleId,
			}).Send(ctx.Context)
			if err != nil {
				return rg, fmt.Errorf("failed to get rate based rule %s: %w", aws.StringValue(id.RuleId), err)
			}
			if v := rule.Rule; v != nil {
				r := resource.New(ctx, resource.WafRegionalRateBasedRule, v.RuleId, v.Name, v)
				rg.AddResource(r)
			}
		}
		if res.NextMarker == nil {
			break
		}
		nextMarker = res.NextMarker
	}

	return rg, nil
}
