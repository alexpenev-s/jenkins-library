import com.sap.piper.Credential
import groovy.transform.Field

@Field String STEP_NAME = getClass().getName()
@Field String METADATA_FILE = 'metadata/kaniko.yaml'

void call(Map parameters = [:]) {
    List credentials = [[type: Credential.FILE, id: 'dockerConfigJsonCredentialsId', env: ['PIPER_dockerConfigJSON']]]
    piperExecuteBin(parameters, STEP_NAME, METADATA_FILE, credentials)
}
