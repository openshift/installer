#!/bin/env groovy

folder("maintenance")

// TAG AWS - Jenkins job to tag aws resources using Grafiti
job("maintenance/tag-aws") {
  description('Tags AWS resources across all AWS regions every 6 hours with an \'expirationDate\' of today, if not already tagged. Changes here will be reverted automatically.')

  label 'worker&&ec2'
  logRotator(30, 100)

  parameters {
    stringParam('START_HOUR', '7', 'Number of hours prior to now to start parsing logs from CloudTrail.')
    stringParam('END_HOUR', '0', 'Number of hours prior to now to stop parsing logs from CloudTrail.')
    stringParam('AWS_REGION', '', 'Optional. Specific AWS region in which to tag resources.')
    stringParam('DATE_VALUE_OVERRIDE', '', 'Optional. YYYY-MM-DD formatted tag value of resources to delete.')
    stringParam('GRAFITI_IMAGE', 'quay.io/coreos/grafiti:64182b19f2c852ab8351ab56bf9533b3b99dfc63', 'Grafiti docker image to use in the job')
  }

  wrappers {
    colorizeOutput()
    timestamps()
    credentialsBinding {
      usernamePassword("QUAY_USERNAME", "QUAY_PASSWD", "quay-robot")
      usernamePassword("AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "jenkins-tectonic-installer")
    }
  }

  triggers {
    cron('H H/6 * * *')
  }

  steps {
    def cmd = """#!/bin/bash -e
date_string=\"\$(date \"+%Y-%m-%d\")\"

if [ -n \"\${DATE_VALUE_OVERRIDE}\" ]; then
    date_string=\"\${DATE_VALUE_OVERRIDE}\"
fi

tag_file=\$(mktemp)
cat <<EOF > \"\$tag_file\"
    [{\"Key\": \"expirationDate\",\"Value\": \"\$date_string\"}]
EOF

private_zones=\$(aws route53 list-hosted-zones | \\
              jq \".HostedZones[] | select(.Config.PrivateZone == true) | .Id\" | \\


for key in \$(cat \"\$tag_file\" | jq \".[].Key\"); do
  for zone in \$private_zones; do
  is_not_tagged=\$(aws route53 list-tags-for-resource \\
                  --resource-type hostedzone \\
                  --resource-id \$(echo \${zone##*/} | \\
                  sed \"s@\\\"@@g\") | \\
                  jq \".ResourceTagSet | select(.Tags[]? | .Key == \$key) | .ResourceId\")
    if [ -z \"\$is_not_tagged\" ]; then
      set +e
      aws route53 change-tags-for-resource \\
      --resource-type hostedzone \
      --add-tags \"\$(cat \"\$tag_file\")\" \\
      --resource-id \"\${zone##*/}\"
      set -e
      echo \"Tagged hosted zone \${zone##*/}\"
    fi
  done
done
rm -f \"\$tag_file\"
    """.stripIndent()

    shell(cmd)
  }

  steps {
    def cmd = """#!/bin/bash
GRF_IMAGE=\${GRAFITI_IMAGE}

cfg_dir=\"\$WORKSPACE\"/config
mkdir -p \"\$cfg_dir\"

regions=( \"us-east-2\" \"us-east-1\" \"us-west-1\" \"us-west-2\" \"ca-central-1\"
\"ap-south-1\" \"ap-northeast-2\" \"ap-southeast-1\" \"ap-southeast-2\"
\"ap-northeast-1\" \"eu-central-1\" \"eu-west-1\" \"eu-west-2\" \"sa-east-1\" )
if [ -n \"\$AWS_REGION\" ]; then
  regions=( \"\$AWS_REGION\" )
fi

date_string='now|strftime(\\\"%Y-%m-%d\\\")'
if [ -n \"\$DATE_VALUE_OVERRIDE\" ]; then
  date_string='\\\"'\"\${DATE_VALUE_OVERRIDE}\"'\\\"'
fi

# Configure grafiti to tag all resources created in the previous 6 hour
# period with an expirationDate of today
cat <<EOF > \"\${cfg_dir}/config.toml\"
endHour = -\${END_HOUR}
startHour = -\${START_HOUR}
includeEvent = false
tagPatterns = [
  \"{expirationDate: \${date_string}}\"
]
EOF

# Exclusion file prevents tagging of resources that already have tags with the key
# \"expirationDate\"
cat <<EOF > \"\${cfg_dir}/filter.json\"
{
  \"TagFilters\": [
    {
      \"Key\": \"expirationDate\",
      \"Values\": []
    }
  ]
}
EOF

for r in \"\${regions[@]}\"; do
  echo \"Tagging AWS region \\\"\$r\\\"\"
  docker run -t --rm \\
    -v \"\$cfg_dir\":/tmp/config:z \\
    -e AWS_ACCESS_KEY_ID=\"\$AWS_ACCESS_KEY_ID\" \\
    -e AWS_SECRET_ACCESS_KEY=\"\$AWS_SECRET_ACCESS_KEY\" \\
    -e AWS_REGION=\"\$r\" \\
    -e CONFIG_FILE=\"/tmp/config/config.toml\" \\
    -e TAG_FILE=\"/tmp/config/filter.json\" \\
    \"\$GRF_IMAGE\" \\
    bash -c \"grafiti -c \\\"\\\$CONFIG_FILE\\\" parse | \\
    grafiti -c \\\"\\\$CONFIG_FILE\\\" filter --ignore-file \\\"\\\$TAG_FILE\\\" | \\
    grafiti -c \\\"\\\$CONFIG_FILE\\\" tag\"
done
    """.stripIndent()
    shell(cmd)
  }

  publishers {
    wsCleanup()
    slackNotifier {
      authTokenCredentialId('tectonic-slack-token')
      customMessage("Jenkins Maintenance: tag-aws")
      includeCustomMessage(true)
      notifyBackToNormal(true)
      notifyFailure(true)
      room('#tectonic-installer-ci')
      teamDomain('coreos')
    }
  }
}


// CLEAN AWS - Jenkins job to clean aws resources using Grafiti
job("maintenance/clean-aws") {
  description('Delete AWS resources tagged with an \'expirationDate\' across all AWS regions every 6 hours. Changes here will be reverted automatically.')

  label 'worker&&ec2'
  logRotator(30, 100)

  parameters {
    stringParam('DATE_VALUE_OVERRIDE', '', 'Optional. YYYY-MM-DD formatted tag value of resources to delete..')
    stringParam('AWS_REGION', '', 'Optional. Specific AWS region to clean.')
    stringParam('GRAFITI_IMAGE', 'quay.io/coreos/grafiti:64182b19f2c852ab8351ab56bf9533b3b99dfc63', 'Grafiti docker image to use in the job')
  }

  wrappers {
    colorizeOutput()
    timestamps()
    credentialsBinding {
      usernamePassword("QUAY_USERNAME", "QUAY_PASSWD", "quay-robot")
      usernamePassword("AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "jenkins-tectonic-installer")
    }
  }

  triggers {
    cron('H H/6 * * *')
  }

  steps {
    def cmd = """#!/bin/bash
GRF_IMAGE=\${GRAFITI_IMAGE}

cfg_dir=\"\${WORKSPACE}/config\"
mkdir -p \"\$cfg_dir\"

regions=( \"us-east-2\" \"us-east-1\" \"us-west-1\" \"us-west-2\" \"ca-central-1\"
\"ap-south-1\" \"ap-northeast-2\" \"ap-southeast-1\" \"ap-southeast-2\"
\"ap-northeast-1\" \"eu-central-1\" \"eu-west-1\" \"eu-west-2\" \"sa-east-1\" )

if [ -n \"\$AWS_REGION\" ]; then
  regions=( \"\$AWS_REGION\" )
fi

cat <<EOF > \"\${cfg_dir}/config.toml\"
maxNumRequestRetries = 11
EOF

# Tag file specifies which resources to delete by a tag of
# expirationDate = {today, yesterday}, or by a single date
# if DATE_VALUE_OVERRIDE is not empty
exp_date=\"\$(date \"+%Y-%m-%d\" -d \"-1 day\")\\\",\\\"\$(date \"+%Y-%-m-%-d\" -d \"-1 day\")\\\",\\\"\$(date +%Y-%m-%d)\\\",\\\"\$(date +%Y-%-m-%-d)\"
if [ -n \"\$DATE_VALUE_OVERRIDE\" ]; then
  exp_date=\"\$DATE_VALUE_OVERRIDE\"
fi

cat <<EOF > \"\${cfg_dir}/delete.json\"
{
  \"TagFilters\": [
    {
      \"Key\": \"expirationDate\",
      \"Values\": [\"\${exp_date:?exp_date cannot be null.}\"]
    }
  ]
}
EOF

echo \"Excluded tags file:\"
cat \"\$cfg_dir\"/delete.json

for r in \"\${regions[@]}\"; do
  echo \"Deleting AWS region \\\"\$r\\\"\"
  docker run -t --rm \\
    -v \"\$WORKSPACE\"/config:/tmp/config:z \\
    -e AWS_ACCESS_KEY_ID=\"\$AWS_ACCESS_KEY_ID\" \\
    -e AWS_SECRET_ACCESS_KEY=\"\$AWS_SECRET_ACCESS_KEY\" \\
    -e AWS_REGION=\"\$r\" \\
    -e CONFIG_FILE=\"/tmp/config/config.toml\" \\
    -e TAG_FILE=\"/tmp/config/delete.json\" \\
    \"\$GRF_IMAGE\" \\
    bash -c \"grafiti -c \\\$CONFIG_FILE -e delete -s --report --all-deps -f \\\$TAG_FILE\"
done
   """.stripIndent()

    shell(cmd)
  }

  publishers {
    wsCleanup()
    slackNotifier {
      authTokenCredentialId('tectonic-slack-token')
      customMessage("Jenkins Maintenance: clean-aws")
      includeCustomMessage(true)
      notifyBackToNormal(true)
      notifyFailure(true)
      room('#tectonic-installer-ci')
      teamDomain('coreos')
    }
  }
}
