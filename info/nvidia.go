package info

import (
	"bytes"
	"context"
	"encoding/xml"
	"github.com/sirupsen/logrus"
	"os/exec"
)

type NvidiaSMILog struct {
	Timestamp     string `xml:"timestamp"`
	DriverVersion string `xml:"driver_version"`
	CudaVersion   string `xml:"cuda_version"`
	AttachedGPUs  int    `xml:"attached_gpus"`
	Gpus          []GPU  `xml:"gpu"`
}

type GPU struct {
	ID                        string                 `xml:"id,attr"`
	ProductName               string                 `xml:"product_name"`
	ProductBrand              string                 `xml:"product_brand"`
	ProductArchitecture       string                 `xml:"product_architecture"`
	DisplayMode               string                 `xml:"display_mode"`
	DisplayActive             string                 `xml:"display_active"`
	PersistenceMode           string                 `xml:"persistence_mode"`
	AddressingMode            string                 `xml:"addressing_mode"`
	MIGMode                   MIGMode                `xml:"mig_mode"`
	MIGDevices                string                 `xml:"mig_devices"`
	AccountingMode            string                 `xml:"accounting_mode"`
	AccountingModeBufferSize  int                    `xml:"accounting_mode_buffer_size"`
	DriverModel               DriverModel            `xml:"driver_model"`
	Serial                    string                 `xml:"serial"`
	UUID                      string                 `xml:"uuid"`
	MinorNumber               string                 `xml:"minor_number"`
	VBIOSVersion              string                 `xml:"vbios_version"`
	MultiGPUBoard             string                 `xml:"multigpu_board"`
	BoardID                   string                 `xml:"board_id"`
	BoardPartNumber           string                 `xml:"board_part_number"`
	GPUPartNumber             string                 `xml:"gpu_part_number"`
	GPUFRUPartNumber          string                 `xml:"gpu_fru_part_number"`
	GPUModuleID               string                 `xml:"gpu_module_id"`
	InforomVersion            InforomVersion         `xml:"inforom_version"`
	InforomBBXFlush           InforomBBXFlush        `xml:"inforom_bbx_flush"`
	GPUOperationMode          GPUOperationMode       `xml:"gpu_operation_mode"`
	GSPFirmwareVersion        string                 `xml:"gsp_firmware_version"`
	GPUVirtualizationMode     GPUVirtualizationMode  `xml:"gpu_virtualization_mode"`
	GPUResetStatus            GPUResetStatus         `xml:"gpu_reset_status"`
	IBMNPU                    IBMNPU                 `xml:"ibmnpu"`
	PCI                       PCI                    `xml:"pci"`
	FanSpeed                  string                 `xml:"fan_speed"`
	PerformanceState          string                 `xml:"performance_state"`
	ClocksEventReasons        ClocksEventReasons     `xml:"clocks_event_reasons"`
	SparseOperationMode       string                 `xml:"sparse_operation_mode"`
	FBMemoryUsage             MemoryUsage            `xml:"fb_memory_usage"`
	BAR1MemoryUsage           MemoryUsage            `xml:"bar1_memory_usage"`
	CCProtectedMemoryUsage    MemoryUsage            `xml:"cc_protected_memory_usage"`
	ComputeMode               string                 `xml:"compute_mode"`
	Utilization               Utilization            `xml:"utilization"`
	EncoderStats              EncoderStats           `xml:"encoder_stats"`
	FBCStats                  FBCStats               `xml:"fbc_stats"`
	ECCMode                   ECCMode                `xml:"ecc_mode"`
	ECCErrors                 ECCErrors              `xml:"ecc_errors"`
	RetiredPages              RetiredPages           `xml:"retired_pages"`
	RemappedRows              string                 `xml:"remapped_rows"`
	Temperature               Temperature            `xml:"temperature"`
	SupportedGPUTargetTemp    SupportedGPUTargetTemp `xml:"supported_gpu_target_temp"`
	GPUPowerReadings          PowerReadings          `xml:"gpu_power_readings"`
	ModulePowerReadings       PowerReadings          `xml:"module_power_readings"`
	Clocks                    Clocks                 `xml:"clocks"`
	ApplicationsClocks        ApplicationsClocks     `xml:"applications_clocks"`
	DefaultApplicationsClocks ApplicationsClocks     `xml:"default_applications_clocks"`
	DeferredClocks            DeferredClocks         `xml:"deferred_clocks"`
	MaxClocks                 MaxClocks              `xml:"max_clocks"`
	MaxCustomerBoostClocks    MaxCustomerBoostClocks `xml:"max_customer_boost_clocks"`
	ClockPolicy               ClockPolicy            `xml:"clock_policy"`
	Voltage                   Voltage                `xml:"voltage"`
	Fabric                    Fabric                 `xml:"fabric"`
	SupportedClocks           []SupportedMemClock    `xml:"supported_clocks>supported_mem_clock"`
	Processes                 []ProcessInfo          `xml:"processes>process_info"`
	AccountedProcesses        string                 `xml:"accounted_processes"`
}

type MIGMode struct {
	CurrentMIG string `xml:"current_mig"`
	PendingMIG string `xml:"pending_mig"`
}

type DriverModel struct {
	CurrentDM string `xml:"current_dm"`
	PendingDM string `xml:"pending_dm"`
}

type InforomVersion struct {
	ImgVersion string `xml:"img_version"`
	OEMObject  string `xml:"oem_object"`
	ECCObject  string `xml:"ecc_object"`
	PWRObject  string `xml:"pwr_object"`
}

type InforomBBXFlush struct {
	LatestTimestamp string `xml:"latest_timestamp"`
	LatestDuration  string `xml:"latest_duration"`
}

type GPUOperationMode struct {
	CurrentGOM string `xml:"current_gom"`
	PendingGOM string `xml:"pending_gom"`
}

type GPUVirtualizationMode struct {
	VirtualizationMode string `xml:"virtualization_mode"`
	HostVGPU           string `xml:"host_vgpu_mode"`
}

type GPUResetStatus struct {
	ResetRequired            string `xml:"reset_required"`
	DrainAndResetRecommended string `xml:"drain_and_reset_recommended"`
}

type IBMNPU struct {
	RelaxedOrderingMode string `xml:"relaxed_ordering_mode"`
}

type PCI struct {
	PCIBus                string         `xml:"pci_bus"`
	PCIDevice             string         `xml:"pci_device"`
	PCIDomain             string         `xml:"pci_domain"`
	PCIDeviceID           string         `xml:"pci_device_id"`
	PCIBusID              string         `xml:"pci_bus_id"`
	PCISubSystemID        string         `xml:"pci_sub_system_id"`
	PCIGPULinkInfo        PCIGPULinkInfo `xml:"pci_gpu_link_info"`
	PCIBridgeChip         PCIBridgeChip  `xml:"pci_bridge_chip"`
	ReplayCounter         string         `xml:"replay_counter"`
	ReplayRolloverCounter string         `xml:"replay_rollover_counter"`
	TXUtil                string         `xml:"tx_util"`
	RXUtil                string         `xml:"rx_util"`
	AtomicCapsInbound     string         `xml:"atomic_caps_inbound"`
	AtomicCapsOutbound    string         `xml:"atomic_caps_outbound"`
}

type PCIGPULinkInfo struct {
	PCIEGen    PCIEGen    `xml:"pcie_gen"`
	LinkWidths LinkWidths `xml:"link_widths"`
}

type PCIEGen struct {
	MaxLinkGen           string `xml:"max_link_gen"`
	CurrentLinkGen       string `xml:"current_link_gen"`
	DeviceCurrentLinkGen string `xml:"device_current_link_gen"`
	MaxDeviceLinkGen     string `xml:"max_device_link_gen"`
	MaxHostLinkGen       string `xml:"max_host_link_gen"`
}

type LinkWidths struct {
	MaxLinkWidth     string `xml:"max_link_width"`
	CurrentLinkWidth string `xml:"current_link_width"`
}

type PCIBridgeChip struct {
	BridgeChipType string `xml:"bridge_chip_type"`
	BridgeChipFW   string `xml:"bridge_chip_fw"`
}

// MemoryUsage represents the <fb_memory_usage> element in the XML.
type MemoryUsage struct {
	Total    string `xml:"total"`
	Reserved string `xml:"reserved"`
	Used     string `xml:"used"`
	Free     string `xml:"free"`
}

type Utilization struct {
	GPUUtil     string `xml:"gpu_util"`
	MemoryUtil  string `xml:"memory_util"`
	EncoderUtil string `xml:"encoder_util"`
	DecoderUtil string `xml:"decoder_util"`
	JpegUtil    string `xml:"jpeg_util"`
	OfaUtil     string `xml:"ofa_util"`
}

type ClocksEventReasons struct {
	ClocksEventReasonGpuIdle                   string `xml:"clocks_event_reason_gpu_idle"`
	ClocksEventReasonApplicationsClocksSetting string `xml:"clocks_event_reason_applications_clocks_setting"`
	ClocksEventReasonSwPowerCap                string `xml:"clocks_event_reason_sw_power_cap"`
	ClocksEventReasonHwSlowdown                string `xml:"clocks_event_reason_hw_slowdown"`
	ClocksEventReasonHwThermalSlowdown         string `xml:"clocks_event_reason_hw_thermal_slowdown"`
	ClocksEventReasonHwPowerBrakeSlowdown      string `xml:"clocks_event_reason_hw_power_brake_slowdown"`
	ClocksEventReasonSyncBoost                 string `xml:"clocks_event_reason_sync_boost"`
	ClocksEventReasonSwThermalSlowdown         string `xml:"clocks_event_reason_sw_thermal_slowdown"`
	ClocksEventReasonDisplayClocksSetting      string `xml:"clocks_event_reason_display_clocks_setting"`
}

// EncoderStats represents the <encoder_stats> element in the XML.
type EncoderStats struct {
	SessionCount   string `xml:"session_count"`
	AverageFPS     string `xml:"average_fps"`
	AverageLatency string `xml:"average_latency"`
}

// FBCStats represents the <fbc_stats> element in the XML.
type FBCStats struct {
	SessionCount   string `xml:"session_count"`
	AverageFPS     string `xml:"average_fps"`
	AverageLatency string `xml:"average_latency"`
}

// ECCMode represents the <ecc_mode> element in the XML.
type ECCMode struct {
	Current string `xml:"current_ecc"`
	Pending string `xml:"pending_ecc"`
}

type ECCErrors struct {
	Volatile  ErrorType `xml:"volatile"`
	Aggregate ErrorType `xml:"aggregate"`
}

type ErrorType struct {
	SramCorrectable   string `xml:"sram_correctable"`
	SramUncorrectable string `xml:"sram_uncorrectable"`
	DramCorrectable   string `xml:"dram_correctable"`
	DramUncorrectable string `xml:"dram_uncorrectable"`
}

type ECCErrorType struct {
	ErrorCount string `xml:"error_count"`
}

type RetiredPages struct {
	SingleBitRetirement BitRetirement `xml:"multiple_single_bit_retirement"`
	DoubleBitRetirement BitRetirement `xml:"double_bit_retirement"`
	PendingBlacklist    string        `xml:"pending_blacklist"`
	PendingRetirement   string        `xml:"pending_retirement"`
}

type BitRetirement struct {
	RetiredCount    string `xml:"retired_count"`
	RetiredPagelist string `xml:"retired_pagelist"`
}

type Temperature struct {
	GpuTemp                string `xml:"gpu_temp"`
	GpuTempTlimit          string `xml:"gpu_temp_tlimit"`
	GpuTempMaxThreshold    string `xml:"gpu_temp_max_threshold"`
	GpuTempSlowThreshold   string `xml:"gpu_temp_slow_threshold"`
	GpuTempMaxGpuThreshold string `xml:"gpu_temp_max_gpu_threshold"`
	GpuTargetTemperature   string `xml:"gpu_target_temperature"`
	MemoryTemp             string `xml:"memory_temp"`
	GpuTempMaxMemThreshold string `xml:"gpu_temp_max_mem_threshold"`
}

type SupportedGPUTargetTemp struct {
	GPUTargetTempMin string `xml:"gpu_target_temp_min"`
	GPUTargetTempMax string `xml:"gpu_target_temp_max"`
}

type PowerReadings struct {
	PowerState          string `xml:"power_state"`
	PowerDraw           string `xml:"power_draw"`
	CurrentPowerLimit   string `xml:"current_power_limit"`
	RequestedPowerLimit string `xml:"requested_power_limit"`
	DefaultPowerLimit   string `xml:"default_power_limit"`
	MinPowerLimit       string `xml:"min_power_limit"`
	MaxPowerLimit       string `xml:"max_power_limit"`
}

// ApplicationsClocks represents the <applications_clocks> element in the XML.
type ApplicationsClocks struct {
	GraphicsClock string `xml:"graphics_clock"`
	MemClock      string `xml:"mem_clock"`
}

// DefaultApplicationsClocks represents the <default_applications_clocks> element in the XML.
type DefaultApplicationsClocks struct {
	GraphicsClock string `xml:"graphics_clock"`
	MemClock      string `xml:"mem_clock"`
}

// DeferredClocks represents the <deferred_clocks> element in the XML.
type DeferredClocks struct {
	GraphicsClock string `xml:"graphics_clock"`
	MemClock      string `xml:"mem_clock"`
}

// MaxClocks represents the <max_clocks> element in the XML.
type MaxClocks struct {
	GraphicsClock string `xml:"graphics_clock"`
	SMClock       string `xml:"sm_clock"`
	MemClock      string `xml:"mem_clock"`
	VideoClock    string `xml:"video_clock"`
}

// MaxCustomerBoostClocks represents the <max_customer_boost_clocks> element in the XML.
type MaxCustomerBoostClocks struct {
	GraphicsClock string `xml:"graphics_clock"`
}

// ClockPolicy represents the <clock_policy> element in the XML.
type ClockPolicy struct {
	AutoBoost        string `xml:"auto_boost"`
	AutoBoostDefault string `xml:"auto_boost_default"`
}

// Voltage represents the <voltage> element in the XML.
type Voltage struct {
	GraphicsVoltage string `xml:"graphics_voltage"`
}

type Fabric struct {
	State  string `xml:"state"`
	Status string `xml:"status"`
}

// Clocks represents the <clocks> element in the XML.
type Clocks struct {
	GraphicsClock string `xml:"graphics_clock"`
	SMClock       string `xml:"sm_clock"`
	MemClock      string `xml:"mem_clock"`
	VideoClock    string `xml:"video_clock"`
}

type SupportedMemClock struct {
	Value                   string   `xml:"value"`
	SupportedGraphicsClocks []string `xml:"supported_graphics_clock"`
}

type ProcessInfo struct {
	GpuInstanceID     string `xml:"gpu_instance_id"`
	ComputeInstanceID string `xml:"compute_instance_id"`
	Pid               string `xml:"pid"`
	Type              string `xml:"type"`
	ProcessName       string `xml:"process_name"`
	UsedMemory        string `xml:"used_memory"`
}

func getNvidiaInfo() ([]byte, error) {
	// 通过执行nvidia-smi -q -xml指令获取返回值，从而能够解析到nvidia显卡的信息
	cmd := exec.CommandContext(context.Background(), "nvidia-smi", "-q", "-x")

	buf := bytes.Buffer{}

	cmd.Stdout = &buf

	if err := cmd.Start(); err != nil {
		logrus.Warnf("start nvidia get info cmd err: %v", err)
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		logrus.Warnf("wait nvidia get info cmd err: %v", err)
		return nil, err
	}

	return buf.Bytes(), nil
}

func Nvidia() (*NvidiaSMILog, error) {
	info, err := getNvidiaInfo()
	if err != nil {
		return nil, err
	}

	//fmt.Printf("%s", string(info))

	result := &NvidiaSMILog{}

	err1 := xml.Unmarshal(info, &result)
	if err1 != nil {
		return nil, err1
	}

	return result, nil
}
