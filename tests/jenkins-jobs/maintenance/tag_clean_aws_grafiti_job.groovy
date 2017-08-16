#!/bin/env groovy​

folder("maintenance")

job("maintenance/tag_clean_aws_resources_grafiti") {
  description("Tag/Clean AWS resources with an \'expirationDate\' of today across all AWS regions every 6 hours.\nThis job is manage by tectonic-installer.\nChanges here will be reverted automatically")
  logRotator(30, 100)
  parameters {
      choiceParam('TAG_CLEAN', ['clean-aws', 'tag-aws'], 'Select \'clean-aws\' to clean the AWS resources or \'tag-aws\' to tag the AWS resources')
      stringParam('GRAFITI_VERSION', '64182b19f2c852ab8351ab56bf9533b3b99dfc63', 'Version of grafiti to run.')
      stringParam('AWS_REGION', '', 'Optional. Specific AWS region to clean.')
      stringParam('DATE_OVERRIDE', '', 'Optional. YYYY-MM-DD formatted tag value of resources to delete.')
      stringParam('START_HOUR', '7', 'Used only when in Tag mode. Number of hours prior to now to start parsing logs from CloudTrail.')
      stringParam('END_HOUR', '0', 'Used only when in Tag mode. Number of hours prior to now to stop parsing logs from CloudTrail.')
      stringParam('SCRIPT_DIR', 'installer/scripts', 'Folder which contains the scripts to run grafiti.')
  }

  label('worker&&ec2')

  scm {
    git {
      remote {
        url('https://github.com/coreos/tectonic-installer')
      }
      branch('origin/master')
    }
  }

  wrappers {
    colorizeOutput()
    timestamps()
    credentialsBinding {
      usernamePassword("QUAY_USERNAME", "QUAY_PASSWD", "quay-robot")
      amazonWebServicesCredentialsBinding {
        accessKeyVariable("AWS_ACCESS_KEY_ID")
        secretKeyVariable("AWS_SECRET_ACCESS_KEY")
        credentialsId("tectonic-jenkins-installer")
      }
    }
  }

  triggers {
    parameterizedCron {
      parameterizedSpecification("""H H/6 * * * % TAG_CLEAN=tag-aws
H H/6 * * * % TAG_CLEAN=clean-aws
    """.stripIndent())
    }
  }

  steps {
    systemGroovyCommand("""
def currentBuild = Thread.currentThread().executable
def description = build.buildVariableResolver.resolve('TAG_CLEAN')
currentBuild.setDescription(description)
  """)

    def cmd = """#!/bin/bash -ex
export DATE_OVERRIDE_FLAG=""
if [ -n "\$DATE_OVERRIDE" ]; then
  export DATE_OVERRIDE_FLAG="--date-override \$DATE_OVERRIDE"
fi

if [ "\$TAG_CLEAN" == "tag-aws" ]; then
  # Tag all Route53 hosted zones
  \$SCRIPT_DIR/maintenance/tag-route53-hosted-zones.sh --force \$DATE_OVERRIDE_FLAG
fi

regions=( "us-east-2" "us-east-1" "us-west-1" "us-west-2" "ca-central-1"
    "ap-south-1" "ap-northeast-2" "ap-southeast-1" "ap-southeast-2"
    "ap-northeast-1" "eu-central-1" "eu-west-1" "eu-west-2" "sa-east-1" )

if [ -n "\$AWS_REGION" ]; then
  regions=( "\$AWS_REGION" )
fi

for region in "\${regions[@]}"; do
  \$SCRIPT_DIR/maintenance/\$TAG_CLEAN.sh \\
    --grafiti-version "\$GRAFITI_VERSION" \\
    --aws-region "\$region" \\
    --workspace-dir "\$WORKSPACE" \\
    --force \\
    \$DATE_OVERRIDE_FLAG
done
    """.stripIndent()
    shell(cmd)
  }

  publishers {
    wsCleanup()
    slackNotifier {
      authTokenCredentialId('tectonic-slack-token')
      customMessage("Tectonic Tag/Clean AWS resources")
      includeCustomMessage(true)
      notifyBackToNormal(true)
      notifyFailure(true)
      room('#tectonic-installer-ci')
      teamDomain('coreos')
    }
  }
}
