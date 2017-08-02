#!/bin/env groovy​

def script_dir = "../../../installer/scripts"
def components = [
  [ name: 'tag-aws',
    description:  "Tag AWS resources with an \'expirationDate\' of today across all AWS regions every 6 hours.",
    extra_commands: """
# Tag all Route53 hosted zones
$script_dir/maintenance/tag-route53-hosted-zones.sh \$date_override_flag

# Tag resources across all AWS regions
    """.stripIndent()
    ],
  [ name: 'clean-aws',
    description: "Delete AWS resources tagged with an \'expirationDate\' across all AWS regions every 6 hours.",
    extra_commands: "" ]
]

folder("maintenance")

components.each {

  def component = it

  job("maintenance/${component.name}") {
    description("${component.description}\nThis job is manage by tectonic-installer.\nChanges here will be reverted automatically")
    logRotator(30, 100)
    label('worker&&ec2')

    wrappers {
      colorizeOutput()
      timestamps()
      credentialsBinding {
        usernamePassword("AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "tectonic-jenkins-installer")
      }
    }

    triggers {
      cron('H H/6 * * *')
    }

    scm {
      git {
        remote {
          url('https://github.com/coreos/tectonic-installer')
        }
        branch('origin/master')
      }
    }

  steps {
      def cmd = """#!/bin/bash -eu
date_override_flag=\"\"
if [ -n \"\$DATE_OVERRIDE\" ]; then
  date_override_flag=\"--date-override \$DATE_OVERRIDE\"
fi

${component.extra_commands}
regions=( \"us-east-2\" \"us-east-1\" \"us-west-1\" \"us-west-2\" \"ca-central-1\"
    \"ap-south-1\" \"ap-northeast-2\" \"ap-southeast-1\" \"ap-southeast-2\"
    \"ap-northeast-1\" \"eu-central-1\" \"eu-west-1\" \"eu-west-2\" \"sa-east-1\" )

if [ -n \"\$AWS_REGION\" ]; then
  regions=( \"\$AWS_REGION\" )
fi

for region in \"\${regions[@]}\"; do
  $script_dir/maintenance/${component.name}.sh \\
    --grafiti-version \"\$GRAFITI_VERSION\" \\
    --aws-region \"\$region\" \\
    --workspace \"\$WORKSPACE\" \\
    --force \\
    \$date_override_flag
done
      """.stripIndent()
      shell(cmd)
    }

    publishers {
      wsCleanup()
      slackNotifier {
        authTokenCredentialId('tectonic-slack-token')
        customMessage("Jenkins Maintenance: ${component.name}")
        includeCustomMessage(true)
        notifyBackToNormal(true)
        notifyFailure(true)
        room('#tectonic-installer-ci')
        teamDomain('coreos')
      }
    }
  }
}
