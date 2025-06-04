package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"text/template"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)


func main() {
	
	// Command line arguments
	templateFile := flag.String("templateFile", "", "The template file to use")
	outputFile := flag.String("outputFile", "", "The output file to write")
	varSource := flag.String("varSource", "none", "Source of template variables (secretsmanager)")
	secretArn := flag.String("secretArn", "", "ARN or Alias of AWS Secrets Manager secret")
	awsRegion := flag.String("region", "", "AWS Region for the service")

	// Parse the command line arguments
	flag.Parse()

	// Validate required arguments
	if *templateFile == "" {
		log.Fatal("templateFile argument is required")
	}

	if *outputFile == "" {
		log.Fatal("outputFile argument is required")
	}

	// Validate variable source if specified
	if *varSource != "secretsmanager" {
		log.Fatal("varSource must be'secretsmanager'")
	}

	// If secrets manager is selected, validate secret ARN and region
	// TODO: Add support for SSM Parameter Store
	if *varSource == "secretsmanager" {
		if *secretArn == "" {
			log.Fatal("secretArn is required when varSource is secretsmanager")
		}
		if *awsRegion == "" {
			log.Fatal("region is required when varSource is secretsmanager")
		}
	}

	fmt.Printf("Template File: %s\n", *templateFile)
	fmt.Printf("Output File: %s\n", *outputFile)
	fmt.Printf("Variable Source: %s\n", *varSource)

	if *varSource == "secretsmanager" {
		fmt.Printf("Secret ARN: %s\n", *secretArn)
		fmt.Printf("AWS Region: %s\n", *awsRegion)
	}

	// parse the template file
	tmpl, err := template.New("config").Option("missingkey=error").ParseFiles(*templateFile)
	if err != nil {
		log.Fatalf("Failed to parse template file: %v", err)
	}

	// Variables to be used in template rendering
	var templateVariables map[string]interface{}

	// if varSource is secretsmanager, get the secret from AWS Secrets Manager
	if *varSource == "secretsmanager" {
		// Load AWS configuration with specified region
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(*awsRegion),
		)
		if err != nil {
			log.Fatalf("Failed to load AWS configuration: %v", err)
		}

		// Create Secrets Manager client
		svc := secretsmanager.NewFromConfig(cfg)

		// Get the secret value
		result, err := svc.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
			SecretId: secretArn,
		})
		if err != nil {
			log.Fatalf("Failed to get secret value: %v", err)
		}

		// Parse the secret JSON into template variables
		if err := json.Unmarshal([]byte(*result.SecretString), &templateVariables); err != nil {
			log.Fatalf("Failed to parse secret JSON: %v", err)
		}
	}

	// Create output file
	outputF, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outputF.Close()
	

	// Execute template with variables
	if err := tmpl.ExecuteTemplate(outputF, *templateFile, templateVariables); err != nil {
		// Check if error is due to missing keys
		if err, ok := err.(template.ExecError); ok {
			log.Fatalf("Template execution failed - missing key in template variables: %v", err)
		}
		log.Fatalf("Failed to execute template: %v", err)
	}

	fmt.Printf("Successfully rendered template to %s\n", *outputFile)
}

